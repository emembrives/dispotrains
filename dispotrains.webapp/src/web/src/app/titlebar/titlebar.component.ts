import { Component, Input } from '@angular/core';

@Component({
  selector: 'titlebar',
  templateUrl: './titlebar.component.html',
  styleUrls: ['./titlebar.component.css']
})
export class TitlebarComponent {
  @Input()
  title: string;

  constructor() { }
}
