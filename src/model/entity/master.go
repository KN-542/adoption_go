package entity

import "api/src/model/ddl"

// m_site
type Site struct {
	ddl.Site
}

// m_schedule_freq_status
type ScheduleFreqStatus struct {
	ddl.ScheduleFreqStatus
}

// m_select_status_event
type SelectStatusEvent struct {
	ddl.SelectStatusEvent
}

// m_assign_rule
type AssignRule struct {
	ddl.AssignRule
}

// m_auto_assign_rule
type AutoAssignRule struct {
	ddl.AutoAssignRule
}
