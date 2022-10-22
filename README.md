# What-Next

`what-next` is a simple tool to help you make better use of the time between meetings.

![a gif demonstrating the basic use of the tool](https://github.com/AP-Hunt/what-next/blob/main/demo.gif?raw=true)

Have you ever wondered what to do in the 20 minutes before your next meeting? `what-next` might be for you. Add your work calendar, and the small jobs you could be doing, and `what-next` will tell you which ones you could get done before your next meeting.

## Installation
Download [the latest release](https://github.com/AP-Hunt/what-next/releases/latest) for your OS from GitHub and install it in your path. It's that easy.

## Adding calendars
`what-next` supports any [ical format](https://icalendar.org/) calendar that it can reach over the internet, or other an ical file on your local machine.

It supports multiple calendars and will show ongoing and upcoming events from all calendars simultaneously.

```sh
$ what-next calendar add work "https://example.org/daily.ical"
```

### Can I use my Google calendar?
You can use your Google calendar! Google helpfully provides [calendars in ical format via a secret link](https://support.google.com/calendar/answer/37648?hl=en#zippy=%2Cget-your-calendar-view-only).

### Can I use my Microsoft Office 365 calendar?
You can use your Micorosft Office 365 calendar! Microsoft has documented [how to get your calendar in ical format](https://support.microsoft.com/en-us/office/introduction-to-publishing-internet-calendars-a25e68d6-695a-41c6-a701-103d44ba151d) on their support site.

## Tracking your todo list
`what-next` has a deliberately simple todo list feature set; it only supports due dates and durations. It's supposed to be quick and unobtrusive to add things, so you can add them before you forget. 

```sh
$ what-next todo add "do the mandatory GDPR training" \
    --due @tomorrow \
    --duration 30m
```
