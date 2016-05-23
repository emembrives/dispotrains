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
        self.statuses = []
    
    def add_status(self, date, desc):
        self.statuses.append((date, desc))

    def finish(self):
        self.statuses.sort(key=lambda x:x[0])

def to_simplified_status(entry):
    try:
        new_entry = {"elevator": entry["elevator"], "state": entry["state"]}
    except TypeError as e:
        print entry
        raise e
    if new_entry["state"].startswith("Hors service"):
        new_entry["state"] = "Hors service"
    elif new_entry["state"].startswith("En travaux"):
        new_entry["state"] = "Hors service"
    elif new_entry["state"].startswith(u"Autre problème"):
        new_entry["state"] = "Hors service"
    new_entry["date"] = datetime.datetime.fromtimestamp(int(entry["lastupdate"]["$date"])/1000)
    return new_entry

simplified_statuses = itertools.imap(to_simplified_status, statuses)

elevators = {}
for status in simplified_statuses:
    elevator = elevators.setdefault(status["elevator"], Elevator(status["elevator"]))
    elevator.add_status(status["date"], status["state"])

days={'Hors service': {'ratp': [], 'sncf': []}, 'Disponible': {'ratp': [], 'sncf': []}}
for elevator in elevators.values():
    elevator.finish()
    current_state=None
    start=0
    last_state=0
    for state in elevator.statuses:
        last_state=state[0]
        if state[1] != "Hors service" and state[1] != "Disponible":
            continue
        if state[1] != current_state:
            if current_state != None:
                delta = state[0]-start
                days[current_state]['sncf' if elevator.name.isdigit() else 'ratp'].append(delta.days * 24 + delta.seconds/3600)
            start = state[0]
            current_state = state[1]
    if start != 0:
        delta = last_state-start
        days[current_state]['sncf' if elevator.name.isdigit() else 'ratp'].append(delta.days * 24 + delta.seconds/3600)

for s, n in days.items():
    for a, ss in n.items():
        for d in ss:
            print "%s,%s,%d" % (s, a, d)
