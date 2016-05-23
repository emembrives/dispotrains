import json
import datetime
import csv

f = open("../../data/statuses.json")

UNKNOWN_STATE = "Information non disponible"
AVAILABLE_STATE = "Disponible"
PLANNED_STATE = "jusqu'au"

class Status(object):
  def __init__(self, date, state):
    self.date = date
    self.state = state

  def date_key(self):
    return self.date

  def is_known(self):
    return self.state != UNKNOWN_STATE

  def is_available(self):
    return self.state == AVAILABLE_STATE

  def is_planned(self):
    return PLANNED_STATE in self.state

class Elevator(object):
  def __init__(self, name):
    self.name = name
    self.statuses = []
    self.unavailable_prob = -1
    self.failure_prob = -1
    self.repair_prob = -1

  def probabilities(self):
    last_status = self.statuses[0]
    total = len(self.statuses)
    unavailable = 0
    failures = 0
    repair = 0
    for status in self.statuses:
      if not status.is_known() or status.is_planned():
        continue
      if not status.is_available():
        unavailable += 1
      if last_status.is_available() and not status.is_available():
        failures += 1
      if not last_status.is_available() and status.is_available():
        repair += 1
      last_status = status
    self.unavailable_prob = float(unavailable)/float(total)
    self.failure_prob = float(failures)/float(total-unavailable) if total-unavailable != 0 else 1.0
    self.repair_prob = float(repair)/float(unavailable) if unavailable != 0 else 1.0

  def add_status(self, date, state):
    self.statuses.append(Status(
        datetime.datetime.fromtimestamp(int(date)/1000),
        state))

  def add_missing_statuses(self):
    self.statuses.sort(key=Status.date_key)
    extra_statuses = []

    current_day = self.statuses[0].date.date()
    num_reports = 0
    current_index = 0
    while current_index < len(self.statuses):
      status = self.statuses[current_index]
      if status.date.date() != current_day:
        if num_reports < 3:
          extra_status = Status(datetime.datetime.combine(
            current_day, datetime.time(23, 0, 0)), UNKNOWN_STATE)
          extra_statuses.extend([extra_status] * (num_reports - 3))
        current_day += datetime.timedelta(1)
        num_reports = 0
      else:
        num_reports += 1
        current_index += 1
    self.statuses.extend(extra_statuses)
    self.statuses.sort(key=Status.date_key)

  def is_sncf(self):
    return self.name.isdigit()


elevators = {}
for l in f.readlines():
  j = json.loads(l.strip())
  if j["elevator"] not in elevators:
    elevators[j["elevator"]] = Elevator(j["elevator"])
  elevators[j["elevator"]].add_status(j["lastupdate"]["$date"], j["state"])

# Length of a status sequence.
SEQUENCE_LENGTH = 30

with open('data/data.csv', 'wb') as csvfile:
  f = csv.writer(csvfile)
  #f.writerow(["name", "sncf"] + ["S-" + str(i) for i in range(SEQUENCE_LENGTH)])
  elevator_index = 0
  for elevator in elevators.values():
    print elevator.name
    elevator_index += 1
    elevator.add_missing_statuses()
    elevator.probabilities()
    for index in range(len(elevator.statuses) - SEQUENCE_LENGTH):
      statuses = elevator.statuses[index : index + SEQUENCE_LENGTH]
      if (not all(map(Status.is_known, statuses))) or any(map(Status.is_planned, statuses)):
        # Skip planned or unknown states
        continue
      f.writerow([elevator_index, int(elevator.is_sncf()), elevator.unavailable_prob, elevator.failure_prob, elevator.repair_prob] + [int(status.is_available()) for status in statuses])
