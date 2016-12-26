import { Component, Input } from '@angular/core';
import { Observable }       from 'rxjs/Observable';

import { Elevator } from '../station';

@Component({
  selector: 'elevator-item',
  templateUrl: './elevator-item.component.html',
  styleUrls: ['./elevator-item.component.css']
})
export class ElevatorItemComponent {
  @Input()
  elevator: Observable<Elevator>;

  constructor() { }
}
