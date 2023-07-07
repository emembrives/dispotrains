export class ElevatorStatistics {
  Mtbf: number;
  Mtbr: number;
  Broken: number;
  Total: number;
  States: ElevatorState[];

  constructor(es: ElevatorStatistics) {
    this.Mtbf = es.Mtbf;
    this.Mtbr = es.Mtbr;
    this.Broken = es.Broken;
    this.Total = es.Total;
    this.States = es.States.map((s) => new ElevatorState(s));
  }
}

export class ElevatorState {
  elevator: string;
  state: string;
  begin: Date;
  end: Date;

  constructor(es: ElevatorState) {
    this.elevator = es.elevator;
    this.state = es.state;
    this.begin = es.begin;
    this.end = es.end;
  }

  isBroken() {
    return this.state !== "Disponible";
  }
}
