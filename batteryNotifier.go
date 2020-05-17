package main

import (
	"github.com/distatus/battery"
	"github.com/mqu/go-notify"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"time"
	"os"
)

func playBeep() {
	f, err := os.Open("resources/.scripts/beep.mp3")
	if err == nil {
		streamer, format, err := mp3.Decode(f)
		if err == nil {
			defer streamer.Close()
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

			done := make(chan bool)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				done <- true
			})))
			<- done
			close(done)
			f.Close()
		}
	}
}

func playInbox() {
	f, err := os.Open("resources/filling-your-inbox.mp3")
	if err == nil {
		streamer, format, err := mp3.Decode(f)
		if err == nil {
			defer streamer.Close()
			speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

			done := make(chan bool)
			speaker.Play(beep.Seq(streamer, beep.Callback(func() {
				done <- true
			})))
			<- done
			close(done)
			f.Close()
		}
	}
}

func notifyFullBattery(notified *bool) {
	if (!*(notified)) {
		*notified = true
		fullBattery := notify.NotificationNew("Battery Full", "", "dialog-information")
		notify.NotificationSetUrgency(fullBattery, notify.NOTIFY_URGENCY_LOW)
		fullBattery.Show()
		playBeep()
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
		playBeep()
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
		playInbox()
		time.Sleep(2 * time.Second)
		notify.NotificationClose(charging)
	} else if (!isCharging && *chargeNotified) {
		*chargeNotified = false
		discharging := notify.NotificationNew("Battery Discharging", "", "dialog-information")
		notify.NotificationSetUrgency(discharging, notify.NOTIFY_URGENCY_NORMAL)
		discharging.Show()
		playInbox()
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
