export class Line {
  network: string;
  id: string;

  constructor(lineData: Object) {
    this.network = lineData["network"];
    this.id = lineData["id"];
  }

  GetName(): string {
    return this.network + " " + this.id;
  }
}

export class Position {
  latitude: number;
  longitude: number;

  constructor(positionData: Object) {
    this.latitude = positionData["latitude"];
    this.longitude = positionData["longitude"];
  }
}

export class Status {
  lastupdate: string;
  state: string;

  constructor(statusData: Object) {
    this.lastupdate = statusData["lastUpdate"];
    this.state = statusData["state"];
  }
}

export class Elevator {
  direction: string;
  id: string;
  situation: string;
  status: Status;

  constructor(elevatorData: Object) {
    this.direction = elevatorData["direction"];
    this.id = elevatorData["id"];
    this.situation = elevatorData["situation"];
    this.status = new Status(elevatorData["status"]);
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
    this.name = stationData["name"];
    this.displayname = stationData["displayname"];
    this.city = stationData["city"];
    this.osmid = stationData["osmid"];
    this.position = new Position(stationData["position"]);

    this.lines = new Array<Line>();
    for (let lineData of stationData["lines"]) {
      this.lines.push(new Line(lineData));
    }

    this.elevators = new Array<Elevator>();
    for (let elevatorData of stationData["elevators"]) {
      this.elevators.push(new Elevator(elevatorData));
    }
  }
}
