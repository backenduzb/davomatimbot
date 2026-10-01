package states

const (
	StateWaitingAdminMenu        = "waiting_admin_menu"
	StateWaitingAdminBroadcast   = "waiting_admin_broadcast"
	StateWaitingAdminClassChoice   = "waiting_admin_class_choice"
	StateWaitingAdminTeacherChoice = "waiting_admin_teacher_choice"
	StateWaitingAbsentTypeChoice = "waiting_absent_type_choice"
	StateWaitingAbsentStudent    = "waiting_absent_student"
	StateWaitingAbsentConfirm    = "waiting_absent_confirm"

	StateWaitingReasonStudent = "waiting_reason_student"
	StateWaitingReasonInput   = "waiting_reason_input"
	StateWaitingReasonConfirm = "waiting_reason_confirm"

	StateWaitingLateStudent = "waiting_late_student"
	StateWaitingLateConfirm = "waiting_late_confirm"

	StateWaitingGuestTeacherChoice = "waiting_guest_teacher_choice"
	StateWaitingGuestCongratsInput = "waiting_guest_congrats_input"
)
