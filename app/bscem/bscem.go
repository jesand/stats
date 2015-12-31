package main

import (
	"encoding/csv"
	"fmt"
	"github.com/docopt/docopt-go"
	"github.com/jesand/stats/dist"
	"github.com/jesand/stats/model"
	"io"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

const (
	USAGE = `bscem - Run Expectation Maximization to infer parameters and truth
        for a binary symmetric channel

Usage:
  bscem <prefs> <qrel> <research_task> <topic> [--baseline] [--pair] [--soft]

Options:
  <prefs>          A CSV file containing channel output
  <qrel>           The QREL containing gold standard assessments
  <research_task>  The research task to assess
  <topic>          The topic to assess
  --baseline       Use the majority vote and const-resp models
  --pair           Use the BSCPair model
  --soft           Use soft assignments during inference
`

	FldTaskStr    = "research_task"
	FldTopicStr   = "topic"
	FldTopicIdStr = "topic_id"
	FldStatusStr  = "assn_status"
	FldWorkerStr  = "worker_id"
	FldDoc1Str    = "left_doc"
	FldDoc2Str    = "right_doc"
	FldVoteStr    = "result"

	ValApproved = "Approved"
	ValLt       = "win"

	InitialNoise = 1e-3
)

func main() {

	// Parse the command line
	args, _ := docopt.Parse(USAGE, nil, true, "1.0", false)
	var (
		file        = args["<prefs>"].(string)
		qrelPath    = args["<qrel>"].(string)
		taskFilter  = args["<research_task>"].(string)
		topicFilter = args["<topic>"].(string)
		baseline    = args["--baseline"].(bool)
		pairModel   = args["--pair"].(bool)
		soft        = args["--soft"].(bool)
		fields      []string
		rows        [][]string
	)

	fmt.Println("Training ", topicFilter, "with baselines?", baseline, "with pair?", pairModel, "with soft?", soft)

	// Initialize
	if f, err := os.Open(file); err != nil {
		fmt.Println("Could not open ", file, ":", err)
		return
	} else {
		r := csv.NewReader(f)

		fmt.Println("Loading", file)
		if fields, err = r.Read(); err != nil {
			fmt.Println("Could not read CSV field names - ", err)
			return
		} else if rows, err = r.ReadAll(); err != nil {
			fmt.Println("Could not read CSV contents - ", err)
			return
		}
		f.Close()
	}

	// Find our fields
	var (
		fldTask    int
		fldTopic   int
		fldTopicId int
		fldStatus  int
		fldWorker  int
		fldDoc1    int
		fldDoc2    int
		fldVote    int

		topicIdFilter string
	)
	for i, field := range fields {
		if field == FldTaskStr {
			fldTask = i
		} else if field == FldTopicStr {
			fldTopic = i
		} else if field == FldTopicIdStr {
			fldTopicId = i
		} else if field == FldStatusStr {
			fldStatus = i
		} else if field == FldWorkerStr {
			fldWorker = i
		} else if field == FldDoc1Str {
			fldDoc1 = i
		} else if field == FldDoc2Str {
			fldDoc2 = i
		} else if field == FldVoteStr {
			fldVote = i
		}
	}
	for _, row := range rows {
		if row[fldTopic] == topicFilter {
			topicIdFilter = row[fldTopicId]
			break
		}
	}

	// Load the QREL
	var qrel = make(map[string]int)
	if f, err := os.Open(qrelPath); err != nil {
		fmt.Println("Could not open ", qrelPath, ":", err)
		return
	} else {
		fmt.Println("Loading", qrelPath)
		r := csv.NewReader(f)
		r.Comma = ' '
		r.FieldsPerRecord = 4
		for {
			record, err := r.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Println("Could not parse", qrelPath, ":", err)
				return
			}
			if record[0] == topicIdFilter {
				qrel[record[2]], _ = strconv.Atoi(record[3])
			}
		}
		f.Close()
		fmt.Println("Found assessments for", len(qrel), "documents")
	}

	// Prepare a model for the data
	fmt.Println("Building model...")
	var (
		mod1           = model.NewMultipleBSCModel()
		mod2           = model.NewMultipleBSCPairModel()
		majority       = make(map[string]map[string]int)
		numPos, numNeg int
	)
	mod1.SoftInputs = soft
	mod2.SoftInputs = soft
	for _, row := range rows {
		if row[fldTask] != taskFilter || row[fldTopic] != topicFilter ||
			row[fldStatus] != ValApproved || row[fldVote] != ValLt {
			continue
		}
		var (
			winner, loser = row[fldDoc1], row[fldDoc2]
			small, big    = winner, loser
			isLt          bool
		)
		if big < small {
			small, big = big, small
		}
		if _, ok := majority[small]; !ok {
			majority[small] = make(map[string]int)
		}
		if small == winner {
			isLt = true
			numPos++
			majority[small][big]++
		} else {
			isLt = false
			numNeg++
			majority[small][big]--
		}
		var (
			question = fmt.Sprintf("%s %s", small, big)
			worker   = row[fldWorker]
		)
		if pairModel {
			if !mod2.HasChannel(question, worker) {
				mod2.AddChannel(question, InitialNoise, worker, InitialNoise)
			}
			mod2.AddObservation(question, question, worker, isLt)
		} else {
			if !mod1.HasChannel(worker) {
				mod1.AddChannel(worker, InitialNoise)
			}
			mod1.AddObservation(question, worker, isLt)
		}
	}
	fmt.Println(numPos, "Pos", numNeg, "Neg")

	if baseline {
		var all1Correct, all0Correct, randCorrect, majCorrect, total float64
		for small, bigs := range majority {
			for big, score := range bigs {
				var (
					rel0   = qrel[small]
					rel1   = qrel[big]
					ltMaj  = score > 0
					ltAll1 = true
					ltAll0 = false
					ltRand = rand.Float64() > 0.5
				)
				if rel0 != rel1 {
					total++
					if ltMaj && rel0 == 1 {
						majCorrect++
					}
					if ltAll1 && rel0 == 1 {
						all1Correct++
					}
					if ltAll0 && rel0 == 1 {
						all0Correct++
					}
					if ltRand && rel0 == 1 {
						randCorrect++
					}
				}
			}
		}

		fmt.Printf("Majority vote accuracy: %d/%d = %f\n",
			int(majCorrect), int(total), majCorrect/total)
		fmt.Printf("All-true vote accuracy: %d/%d = %f\n",
			int(all1Correct), int(total), all1Correct/total)
		fmt.Printf("All-false vote accuracy: %d/%d = %f\n",
			int(all0Correct), int(total), all0Correct/total)
		fmt.Printf("Random vote accuracy: %d/%d = %f\n",
			int(randCorrect), int(total), randCorrect/total)
	} else if pairModel {

		// Define an evaluation method
		eval := func(mod *model.MultipleBSCPairModel, round int, stage string) {
			var correct, total float64
			for doc, value := range mod.Inputs {
				var (
					docs = strings.Split(doc, " ")
					rel0 = qrel[docs[0]]
					rel1 = qrel[docs[1]]
					lt   = dist.BooleanSpace.BoolValue(value.Outcome())
				)
				if rel0 != rel1 {
					total++
					if lt && rel0 == 1 {
						correct++
					}
				}
			}

			if round == 0 {
				fmt.Printf("%s score: %f accuracy: %d/%d = %f\n", stage,
					mod.Score(), int(correct), int(total), correct/total)
			} else {
				fmt.Printf("Round %d %s score: %f accuracy: %d/%d = %f\n", round, stage,
					mod.Score(), int(correct), int(total), correct/total)
			}
		}

		// Run expectation maximization on the model
		// mod2.UpdateBeta1 = false
		// mod2.UpdateBeta2 = false
		mod2.EM(0, 1e-3, eval)
	} else {

		// Define an evaluation method
		eval := func(mod *model.MultipleBSCModel, round int, stage string) {
			var correct, total float64
			for doc, value := range mod.Inputs {
				var (
					docs = strings.Split(doc, " ")
					rel0 = qrel[docs[0]]
					rel1 = qrel[docs[1]]
					lt   = dist.BooleanSpace.BoolValue(value.Outcome())
				)
				if rel0 != rel1 {
					total++
					if lt && rel0 == 1 {
						correct++
					}
				}
			}

			if round == 0 {
				fmt.Printf("%s score: %f accuracy: %d/%d = %f\n", stage,
					mod.Score(), int(correct), int(total), correct/total)
			} else {
				fmt.Printf("Round %d %s score: %f accuracy: %d/%d = %f\n", round, stage,
					mod.Score(), int(correct), int(total), correct/total)
			}
		}

		// Run expectation maximization on the model
		mod1.EM(0, 1e-3, eval)
	}
}
