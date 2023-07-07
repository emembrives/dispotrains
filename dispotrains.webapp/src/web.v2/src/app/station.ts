export class NetworkStats {
  Good!: number;
  Bad!: number;
  LongBad!: number;

  constructor(ns: NetworkStats) {
    this.Good = ns.Good;
    this.Bad = ns.Bad;
    this.LongBad = ns.LongBad;
  }

  Total(): number {
    return this.Good + this.Bad;
  }

  PercentBad(): number {
    return (this.Bad * 100) / (this.Total());
  }

  PercentLongBad(): number {
    return (this.LongBad * 100) / (this.Total());
  }
}

export class Line {
  network!: string;
  id!: string;

  constructor(line: Line) {
    if (line.network === "LocalTrain") {
      this.network = "Transilien";
    } else if (line.network === "RapidTransit") {
      this.network = "RER";
    } else if (line.network === "RailShuttle") {
      this.network = "Navette";
    } else if (line.network === "Metro") {
      this.network = "Métro";
    } else {
      this.network = line.network;
    }
    this.id = line.id;
  }

  GetName(): string {
    return this.network + ' ' + this.id;
  }
}

export class Position {
  latitude!: number;
  longitude!: number;

  constructor(p: Position) {
    this.latitude = p.latitude;
    this.longitude = p.longitude;
  }
}

export class Status {
  lastupdate!: string;
  state!: string;
  forecast!: string | undefined;

  constructor(status: Status | undefined) {
    if (status !== undefined && status !== null) {
      this.lastupdate = status.lastupdate ?? '';
      this.state = status.state ?? 'Information non disponible';
      this.forecast = status.forecast;
    } else {
      this.lastupdate = '';
      this.state = 'Information non disponible';
      this.forecast = undefined;
    }
  }
}

export class Elevator {
  direction!: string;
  id!: string;
  situation!: string;
  status!: Status;

  constructor(elevator: Elevator) {
    this.direction = elevator.direction;
    this.id = elevator.id;
    this.situation = elevator.situation;
    this.status = new Status(elevator.status);
  }

  available(): boolean {
    return this.status.state === 'Disponible';
  }

  isBroken(): boolean {
    return !this.available();
  }

  hasForecast(): boolean {
    return this.status.forecast !== null;
  }
}

export class Station {
  lines!: Line[];
  elevators!: Elevator[];
  displayname!: string;
  name!: string;
  city: string | undefined;
  position: Position | undefined;
  osmid: string | undefined;

  constructor(station: Station) {
    this.lines = station.lines.map(item => new Line(item));
    this.elevators = station.elevators.map(item => new Elevator(item));
    this.displayname = station.displayname;
    this.name = station.name;
    this.city = station.city;
    this.position = station.position !== undefined ? new Position(station.position) : undefined;
    this.osmid = station.osmid;
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
