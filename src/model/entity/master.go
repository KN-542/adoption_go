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

// m_document_rule
type DocumentRule struct {
	ddl.DocumentRule
}

// m_occupation
type Occupation struct {
	ddl.Occupation
}

// m_interview_processing
type Processing struct {
	ddl.Processing
}
