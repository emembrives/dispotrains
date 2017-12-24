export class ElevatorStats {
  elevator: string;
  state: string;
  start: Date;
  end: Date;

  constructor(elevatorData: Object) {
    this.elevator = elevatorData['elevator'];
    this.state = elevatorData['state'];
    this.start = new Date(elevatorData['startdate']);
    this.end = new Date(elevatorData['enddate']);
  }
}
