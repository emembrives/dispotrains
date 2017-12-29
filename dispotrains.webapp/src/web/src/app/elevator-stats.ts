export class ElevatorStatistics {
  mtbf: number;
  mtbr: number;
  broken: number;
  total: number;
  states: ElevatorState[];

  constructor(data: Object) {
    this.mtbf = data['Mtbf'];
    this.mtbr = data['Mtbr'];
    this.broken = data['Broken'];
    this.total = data['Total'];
    this.states = new Array<ElevatorState>();
    for (let dataState of data['States']) {
      this.states.push(new ElevatorState(dataState));
    }
  }
}

export class ElevatorState {
  elevator: string;
  state: string;
  start: Date;
  end: Date;

  constructor(elevatorData: Object) {
    this.elevator = elevatorData['Elevator'];
    this.state = elevatorData['State'];
    this.start = new Date(elevatorData['Begin']);
    this.end = new Date(elevatorData['End']);
  }
}
