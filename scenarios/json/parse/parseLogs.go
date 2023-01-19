package scenjsonparse

import (
	"errors"
	"fmt"

	mj "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/json/model"
	oj "github.com/multiversx/mx-chain-vm-v1_2-go/scenarios/orderedjson"
)

func (p *Parser) processLogList(logsRaw oj.OJsonObject) ([]*mj.LogEntry, error) {
	logList, isList := logsRaw.(*oj.OJsonList)
	if !isList {
		return nil, errors.New("unmarshalled logs list is not a list")
	}
	var logEntries []*mj.LogEntry
	var err error
	for _, logRaw := range logList.AsList() {
		logMap, isMap := logRaw.(*oj.OJsonMap)
		if !isMap {
			return nil, errors.New("unmarshalled log entry is not a map")
		}
		logEntry := mj.LogEntry{}
		for _, kvp := range logMap.OrderedKV {
			switch kvp.Key {
			case "address":
				logEntry.Address, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log address: %w", err)
				}
			case "identifier":
				logEntry.Identifier, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log identifier: %w", err)
				}
			case "topics":
				logEntry.Topics, err = p.parseCheckBytesList(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log entry topics: %w", err)
				}
			case "data":
				logEntry.Data, err = p.parseCheckBytes(kvp.Value)
				if err != nil {
					return nil, fmt.Errorf("invalid log data: %w", err)
				}
			default:
				return nil, fmt.Errorf("unknown log field: %s", kvp.Key)
			}
		}
		logEntries = append(logEntries, &logEntry)
	}

	return logEntries, nil
}
