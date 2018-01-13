export class NetworkStats {
  good: number;
  bad: number;
  longBad: number;

  constructor(data: Object) {
    this.good = data['Good'];
    this.bad = data['Bad'];
    this.longBad = data['LongBad'];
  }

  Total(): number {
    return this.good + this.bad;
  }

  PercentBad(): number {
    return (this.bad * 100) / (this.Total());
  }

  PercentLongBad(): number {
    return (this.longBad * 100) / (this.Total());
  }
}

export class Line {
  network: string;
  id: string;

  constructor(lineData: Object) {
    this.network = lineData['network'];
    this.id = lineData['id'];
  }

  GetName(): string {
    return this.network + ' ' + this.id;
  }
}

export class Position {
  latitude: number;
  longitude: number;

  constructor(positionData: Object) {
    this.latitude = positionData['latitude'];
    this.longitude = positionData['longitude'];
  }
}

export class Status {
  lastupdate: string;
  state: string;
  forecast: string;

  constructor(statusData: Object) {
    if (statusData !== undefined && statusData !== null) {
      this.lastupdate = statusData['lastupdate'];
      this.state = statusData['state'];
      this.forecast = statusData['forecast'];
    } else {
      this.lastupdate = '';
      this.state = 'Information non disponible';
      this.forecast = undefined;
    }
  }
}

export class Elevator {
  direction: string;
  id: string;
  situation: string;
  status: Status;

  constructor(elevatorData: Object) {
    this.direction = elevatorData['direction'];
    this.id = elevatorData['id'];
    this.situation = elevatorData['situation'];
    this.status = new Status(elevatorData['status']);
  }

  available(): boolean {
    return this.status.state === 'Disponible';
  }
}

export class Station {
  lines: Line[];
  elevators: Elevator[];
  displayname: string;
  name: string;
  city: string;
  position: Position;
  osmid: string;

  constructor(stationData: Object) {
    this.name = stationData['name'];
    this.displayname = stationData['displayname'];
    this.city = stationData['city'];
    this.osmid = stationData['osmid'];
    if (stationData['position'] !== undefined) {
      this.position = new Position(stationData['position']);
    }

    this.lines = new Array<Line>();
    for (let lineData of stationData['lines']) {
      this.lines.push(new Line(lineData));
    }

    this.elevators = new Array<Elevator>();
    for (let elevatorData of stationData['elevators']) {
      this.elevators.push(new Elevator(elevatorData));
    }
  }

  available(): boolean {
    for (let elevator of this.elevators) {
      if (!elevator.available()) {
        return false;
      }
    }
    return true;
  }
}
