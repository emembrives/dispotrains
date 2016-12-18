export class Line {
  network: string;
  id: string;
}

export class Position {
  latitude: number;
  longitude: number;
}

export class Status {
  lastupdate: string;
  state: string;
}

export class Elevator {
  direction: string;
  id: string;
  situation: string;
  status: Status;
}

export class Station {
  lines: Line[];
  elevators: Elevator[];
  displayname: string;
  name: string;
  city: string;
  position: Position;
  osmid: string;
}
