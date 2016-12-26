import { Injectable }     from '@angular/core';
import { Http, Response } from '@angular/http';
import { Station }           from './station';
import { Observable }     from 'rxjs/Observable';
import { combineLatest } from 'rxjs/observable/combineLatest';

import { SorterUtils } from './sorting';

@Injectable()
export class StationService {
  private stationsUrl = 'http://dispotrains.membrives.fr/app/GetStations/';

  constructor(private http: Http) { }

  getStations(): Observable<Station[]> {
    return this.http.get(this.stationsUrl)
      .map(this.extractData)
      .catch(this.handleError);
  }

  getStation(nameObservable: Observable<string>): Observable<Station> {
    return combineLatest(this.getStations(), nameObservable, (stations: Station[], name: string) => {
      for (let station of stations) {
        if (station.name == name) {
          return station;
        }
      }
    });
  }

  private extractData(res: Response): Station[] {
    let body = res.json();
    if (!body) {
      return new Array<Station>();
    }
    let stations: Station[] = new Array<Station>();
    for (let stationData of body) {
      stations.push(new Station(stationData));
    }
    return stations;
  }

  private handleError(error: Response | any) {
    // In a real world app, we might use a remote logging infrastructure
    let errMsg: string;
    if (error instanceof Response) {
      const body = error.json() || '';
      const err = body.error || JSON.stringify(body);
      errMsg = `${error.status} - ${error.statusText || ''} ${err}`;
    } else {
      errMsg = error.message ? error.message : error.toString();
    }
    console.error("Error while retrieving stations: " + errMsg);
    return Observable.throw(errMsg);
  }
}
