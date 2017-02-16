The parser extracts some data from the "Activiteitengegevens ontwikkelingssamenwerking"
found at

	https://www.rijksoverheid.nl/opendata/ontwikkelingssamenwerking

and produces a limited CSV format by mapping activities per transaction.

```
go run main.go < IATIACTIVITIES19972007.xml > IATIACTIVITIES19972007.csv
```
