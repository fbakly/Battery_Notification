package main

import (
	"github.com/distatus/battery"
	"github.com/mqu/go-notify"
	"time"
)

func notifyFullBattery(notified *bool) {
	if (!*(notified)) {
		*notified = true
		fullBattery := notify.NotificationNew("Battery Full", "", "dialog-information")
		notify.NotificationSetUrgency(fullBattery, notify.NOTIFY_URGENCY_LOW)
		fullBattery.Show()
		time.Sleep(2 * time.Second)
		notify.NotificationClose(fullBattery)
	}
}

func notifyLowBattery(notified *bool) {
	if (!(*notified)) {
		*notified = true
		batteryLow := notify.NotificationNew("Battery Low", "Charge Now!!!", "dialog-information")
		notify.NotificationSetUrgency(batteryLow, notify.NOTIFY_URGENCY_CRITICAL)
		batteryLow.Show()
		time.Sleep(2 * time.Second)
		notify.NotificationClose(batteryLow)
	}
}

func notifyCharge(isCharging bool, chargeNotified *bool) {
	if (isCharging && !(*chargeNotified)) {
		*chargeNotified = true
		charging := notify.NotificationNew("Charging battery", "", "dialog-information")
		notify.NotificationSetUrgency(charging, notify.NOTIFY_URGENCY_NORMAL)
		charging.Show()
		time.Sleep(2 * time.Second)
		notify.NotificationClose(charging)
	} else if (!isCharging && *chargeNotified) {
		*chargeNotified = false
		discharging := notify.NotificationNew("Battery Discharging", "", "dialog-information")
		notify.NotificationSetUrgency(discharging, notify.NOTIFY_URGENCY_NORMAL)
		discharging.Show()
		time.Sleep(2 * time.Second)
		notify.NotificationClose(discharging)
	}
}

func getBatteryState(batteries []*battery.Battery, isCharging *bool) {
	length := len(batteries)

	*isCharging = false
	for index := 0; index < length; index++ {
		if (batteries[index].State.String() == "Charging") {
			*isCharging = true
		}
	}
}

func main() {
	chargeNotified := false
	isCharging := false
	notified := false
	notify.Init("Battery Percentage Notifier")
	for  {
		percentage := 0
		batteries, err := battery.GetAll()

		if err != nil {
			return
		}

		getBatteryState(batteries, &isCharging)

		for _, battery := range batteries {
			percentage += int((battery.Current / battery.Full) * 100)
		}

		percentage /= 2

		if (percentage == 100) {
			notifyFullBattery(&notified)
		} else if (percentage <= 20) {
			notifyLowBattery(&notified)
		} else {
			notified = false
		}
	
		notifyCharge(isCharging, &chargeNotified)
	}
}
