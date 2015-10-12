// Copyright 2015, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package tabletservermock provides mock interfaces for tabletserver.
package tabletservermock

import (
	"github.com/youtube/vitess/go/vt/dbconfigs"
	"github.com/youtube/vitess/go/vt/mysqlctl"
	pb "github.com/youtube/vitess/go/vt/proto/query"
	"github.com/youtube/vitess/go/vt/proto/topodata"
	"github.com/youtube/vitess/go/vt/tabletserver"
	"github.com/youtube/vitess/go/vt/tabletserver/queryservice"
)

// BroadcastData is used by the mock Controller to send data
// so the tests can check what was sent.
type BroadcastData struct {
	// TERTimestamp stores the last broadcast timestamp
	TERTimestamp int64

	// RealtimeStats stores the last broadcast stats
	RealtimeStats pb.RealtimeStats
}

// Controller is a mock tabletserver.Controller
type Controller struct {
	// CurrentTarget stores the last known target
	CurrentTarget pb.Target

	// QueryServiceEnabled is a state variable
	QueryServiceEnabled bool

	// InitDBConfigError is the return value for InitDBConfig
	InitDBConfigError error

	// SetServingTypeError is the return value for SetServingType
	SetServingTypeError error

	// IsHealthy is the return value for IsHealthy
	IsHealthyError error

	// ReloadSchemaCount counts how many times ReloadSchema was called
	ReloadSchemaCount int

	// BroadcastData is a channel where we send BroadcastHealth data
	BroadcastData chan *BroadcastData
}

// NewController returns a mock of tabletserver.Controller
func NewController() *Controller {
	return &Controller{
		QueryServiceEnabled: false,
		InitDBConfigError:   nil,
		IsHealthyError:      nil,
		ReloadSchemaCount:   0,
		BroadcastData:       make(chan *BroadcastData, 10),
	}
}

// Register is part of the tabletserver.Controller interface
func (tqsc *Controller) Register() {
}

// AddStatusPart is part of the tabletserver.Controller interface
func (tqsc *Controller) AddStatusPart() {
}

// InitDBConfig is part of the tabletserver.Controller interface
func (tqsc *Controller) InitDBConfig(target *pb.Target, dbConfigs *dbconfigs.DBConfigs, schemaOverrides []tabletserver.SchemaOverride, mysqld mysqlctl.MysqlDaemon) error {
	if tqsc.InitDBConfigError == nil {
		tqsc.CurrentTarget = *target
		tqsc.QueryServiceEnabled = true
	} else {
		tqsc.QueryServiceEnabled = false
	}
	return tqsc.InitDBConfigError
}

// SetServingType is part of the tabletserver.Controller interface
func (tqsc *Controller) SetServingType(tabletType topodata.TabletType, serving bool) error {
	if tqsc.SetServingTypeError == nil {
		tqsc.CurrentTarget.TabletType = tabletType
		tqsc.QueryServiceEnabled = serving
	}
	return tqsc.SetServingTypeError
}

// StartService is part of the tabletserver.Controller interface
func (tqsc *Controller) StartService(*pb.Target, *dbconfigs.DBConfigs, []tabletserver.SchemaOverride, mysqlctl.MysqlDaemon) error {
	tqsc.QueryServiceEnabled = true
	return nil
}

// StopService is part of the tabletserver.Controller interface
func (tqsc *Controller) StopService() {
	tqsc.QueryServiceEnabled = false
}

// IsServing is part of the tabletserver.Controller interface
func (tqsc *Controller) IsServing() bool {
	return tqsc.QueryServiceEnabled
}

// IsHealthy is part of the tabletserver.Controller interface
func (tqsc *Controller) IsHealthy() error {
	return tqsc.IsHealthyError
}

// ReloadSchema is part of the tabletserver.Controller interface
func (tqsc *Controller) ReloadSchema() {
	tqsc.ReloadSchemaCount++
}

//ClearQueryPlanCache is part of the tabletserver.Controller interface
func (tqsc *Controller) ClearQueryPlanCache() {
}

// RegisterQueryRuleSource is part of the tabletserver.Controller interface
func (tqsc *Controller) RegisterQueryRuleSource(ruleSource string) {
}

// UnRegisterQueryRuleSource is part of the tabletserver.Controller interface
func (tqsc *Controller) UnRegisterQueryRuleSource(ruleSource string) {
}

// SetQueryRules is part of the tabletserver.Controller interface
func (tqsc *Controller) SetQueryRules(ruleSource string, qrs *tabletserver.QueryRules) error {
	return nil
}

// QueryService is part of the tabletserver.Controller interface
func (tqsc *Controller) QueryService() queryservice.QueryService {
	return nil
}

// BroadcastHealth is part of the tabletserver.Controller interface
func (tqsc *Controller) BroadcastHealth(terTimestamp int64, stats *pb.RealtimeStats) {
	tqsc.BroadcastData <- &BroadcastData{
		TERTimestamp:  terTimestamp,
		RealtimeStats: *stats,
	}
}