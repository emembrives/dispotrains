#!/usr/bin/python
# vim: set fileencoding=utf-8

import json
import datetime
import itertools

statuses = []
f = open('statuses-2014.json')
for line in f.readlines():
    statuses.append(json.loads(line.strip()))
f.close()

class DataByAgency(object):
    def __init__(self):
        self.sncf = 0
        self.ratp = 0

class Elevator(object):
    def __init__(self, name):
        self.name = name
        self.statuses = {}
    
    def add_status(self, date, desc):
        day = datetime.date(date.year, date.month, date.day)
        status_by_day = self.statuses.setdefault(day, [])
        status_by_day.append((date, desc))

    def ensure_day(self, day):
        status_by_day = self.statuses.setdefault(day, [])

def filter_by_date(start_date, end_date):
    start = datetime.datetime.strptime(start_date, "%Y-%m-%d")
    end = datetime.datetime.strptime(end_date, "%Y-%m-%d")

    def filter_by_provided_date(entry):
        entry_date = datetime.datetime.fromtimestamp(int(entry["lastupdate"]["$date"])/1000)
        return start <= entry_date and entry_date <= end
    return filter_by_provided_date

def to_simplified_status(entry):
    new_entry = {"elevator": entry["elevator"], "state": entry["state"]}
    if new_entry["state"].startswith("Hors service"):
        new_entry["state"] = "Hors service"
    elif new_entry["state"].startswith("En travaux"):
        new_entry["state"] = "Hors service"
    elif new_entry["state"].startswith(u"Autre problème"):
        new_entry["state"] = "Hors service"
    new_entry["date"] = datetime.datetime.fromtimestamp(int(entry["lastupdate"]["$date"])/1000)
    return new_entry

filtered_statuses = itertools.ifilter(filter_by_date("2014-01-01", "2014-07-01"), statuses)
simplified_statuses = itertools.imap(to_simplified_status, filtered_statuses)

elevators = {}
for status in simplified_statuses:
    elevator = elevators.setdefault(status["elevator"], Elevator(status["elevator"]))
    elevator.add_status(status["date"], status["state"])

def ensure_dates_for_elevators(elevators):
    day = datetime.date(2014, 1, 1)
    while day < datetime.date(2014, 7, 1):
        for elevator in elevators.values():
            elevator.ensure_day(day)
        day = day + datetime.timedelta(days=1)
    return elevators

elevators = ensure_dates_for_elevators(elevators)

states_per_day_ratp = {}
states_per_day_sncf = {}
for elevator in elevators.values():
    for day, states in elevator.statuses.items():
        if elevator.name.isdigit():
            states_per_day_sncf[len(states)] = states_per_day_sncf.get(len(states), 0) + 1
        else:
            states_per_day_ratp[len(states)] = states_per_day_ratp.get(len(states), 0) + 1

for i in range(0, max(max(states_per_day_ratp.keys()), max(states_per_day_sncf.keys()))+1):
    print "%s, %s, %s" % (i, states_per_day_sncf.get(i, 0), states_per_day_ratp.get(i, 0))

