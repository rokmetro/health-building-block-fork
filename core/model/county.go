/*
 *   Copyright (c) 2020 Board of Trustees of the University of Illinois.
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package model

//County represents county entity
type County struct {
	ID            string
	Name          string
	StateProvince string
	Country       string

	Guidelines     []Guideline
	Locations      []Location
	CountyStatuses []CountyStatus
}

//CountyStatus represents county status entity
type CountyStatus struct {
	ID          string
	Name        string //red, green, yellow...
	Description string

	County County
}
