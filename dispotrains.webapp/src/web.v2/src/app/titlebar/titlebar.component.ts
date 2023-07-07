import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-titlebar',
  templateUrl: './titlebar.component.html',
  styleUrls: ['./titlebar.component.css']
})
export class TitlebarComponent {
  @Input()
  line: string | undefined;
  @Input()
  station: string | undefined;
  @Input()
  elevator: string | undefined;
}
