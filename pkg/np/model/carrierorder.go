package model

type CarrierOrder struct {
	Delay         int
	UID           int
	OrderTypeID   CarrierOrderType
	NumberOfShips int
}

type CarrierOrderType int

func (cot CarrierOrderType) String() string {
	switch cot {
	case 0:
		return "Do Nothing"
	case 1:
		return "Collect All"
	case 2:
		return "Drop All"
	case 3:
		return "Collect"
	case 4:
		return "Drop"
	case 5:
		return "Drop All But"
	case 6:
		return "Garrison Star"
	default:
		return "Unknown Carrier Order"
	}
}
