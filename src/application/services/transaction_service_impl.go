package services

import (
	"bytes"
	"encoding/csv"
	"html/template"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/DiegoSan99/transaction-processor/src/domain/dto"
	"github.com/DiegoSan99/transaction-processor/src/domain/entity"
	"github.com/aws/aws-lambda-go/events"
)

func (ts *TransactionService) ProcessTransactions(event events.S3Event) error {
	object, err := ts.s3Client.Download(event.Records[0].S3.Bucket.Name, event.Records[0].S3.Object.Key)
	if err != nil {
		ts.log.Errorw("Error downloading file from S3", "error", err)
		return err
	}
	transactions, err := ts.ParseRecords(object)
	if err != nil {
		return err
	}
	ts.log.Infow("Transactions parsed", "transactions", transactions)

	err = ts.SaveTransactions(transactions)
	if err != nil {
		return err
	}
	tr, err := ts.transactionRepo.GetByUserID(1)
	if err != nil {
		ts.log.Errorw("Error getting transactions", "error", err)
		return err
	}
	ts.log.Infow("Transactions retrieved", "transactions", tr)

	report := ts.GenerateReport(tr)
	ts.log.Infow("Report generated", "report", report)

	emailBody := ts.generateEmailBody(*report)
	ts.log.Infow("Email body generated", "emailBody", emailBody)

	imagePath, err := filepath.Abs("stori.png")
	if err != nil {
		ts.log.Errorw("Error getting image path", "error", err)
	}

	err = ts.emailClient.SendEmail(tr[0].Account.Email, "Transaction Report", emailBody, imagePath)
	if err != nil {
		ts.log.Errorw("Error sending email", "error", err)
		return err
	}

	return nil
}

func (ts *TransactionService) ParseRecords(data []byte) ([]dto.TransactionCSV, error) {
	r := csv.NewReader(bytes.NewReader(data))

	_, err := r.Read()
	if err != nil {
		ts.log.Errorw("Error reading CSV header", "error", err)
		return nil, err
	}

	var records []dto.TransactionCSV

	for {
		record, err := r.Read()
		if err != nil {
			break
		}

		parsedDate, err := time.Parse("1/2", record[1])
		if err != nil {
			ts.log.Errorw("Error parsing date", "error", err)
			return nil, err
		}
		currentYear := time.Now().Year()
		parsedDate = time.Date(currentYear, parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())

		transactionValue, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			ts.log.Errorw("Error parsing transaction value", "error", err)
			return nil, err
		}
		trType := record[2][0]

		records = append(records, dto.TransactionCSV{
			Id:          record[0],
			Date:        parsedDate,
			Transaction: transactionValue,
			AccountId:   record[3],
			Type:        string(trType),
		})
	}
	return records, nil
}

func (ts *TransactionService) SaveTransactions(transactions []dto.TransactionCSV) error {
	for _, transactionCSV := range transactions {

		trEntity := entity.Transaction{
			AccountID:  1,
			Type:       transactionCSV.Type,
			IntendedAt: transactionCSV.Date,
			Amount:     transactionCSV.Transaction,
		}

		err := ts.transactionRepo.Create(&trEntity)
		if err != nil {
			ts.log.Errorw("Error saving transaction", "error", err)
			return err
		}
	}
	return nil
}

func (ts *TransactionService) GenerateReport(transactions []entity.Transaction) *dto.Report {
	report := dto.Report{
		TransactionsCount: make(map[string]int),
		AverageDebit:      make(map[string]float64),
		AverageCredit:     make(map[string]float64),
	}

	// solo temps
	debitSums := make(map[string]float64)
	creditSums := make(map[string]float64)
	debitCounts := make(map[string]int)
	creditCounts := make(map[string]int)

	for _, t := range transactions {
		month := t.IntendedAt.Format("January")
		amount := t.Amount
		report.TotalBalance += amount

		if t.Type == "-" {
			debitSums[month] += amount
			debitCounts[month]++
		} else {
			creditSums[month] += amount
			creditCounts[month]++
		}

		report.TransactionsCount[month]++
	}

	for month, sum := range debitSums {
		report.AverageDebit[month] = sum / float64(debitCounts[month])
	}
	for month, sum := range creditSums {
		report.AverageCredit[month] = sum / float64(creditCounts[month])
	}

	return &report
}

func (ts *TransactionService) generateEmailBody(report dto.Report) string {
	tmpl := template.Must(template.New("report").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
	  <title>Transaction Report</title>
	</head>
	<body>
	  <h1>Transaction Summary</h1>
	  <p>Total balance: {{ .TotalBalance }}</p>
	  <h2>Transactions Count</h2>
	  {{ range $month, $count := .TransactionsCount }}
	    <p>{{ $month }}: {{ $count }}</p>
	  {{ end }}
	  <h2>Average Debit</h2>
	  {{ range $month, $amount := .AverageDebit }}
	    <p>{{ $month }}: {{ printf "%.2f" $amount }}</p>
	  {{ end }}
	  <h2>Average Credit</h2>
	  {{ range $month, $amount := .AverageCredit }}
	    <p>{{ $month }}: {{ printf "%.2f" $amount }}</p>
	  {{ end }}
	</body>
	</html>
	`))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, report); err != nil {
		log.Fatal(err)
	}
	return buf.String()
}
