// Copyright 2023 The chessBox Crew
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Пакет usfen реализует формат записи шахматной диаграммы UsFEN.
//
// UsFEN (URL-safe FEN) - это FEN, в которой вместо "/" - "~",
// а вместо пробела - "+".
package usfen

import "strings"

// FromFen конвертирует FEN в UsFEN (без валидации).
func FromFen(fen string) string {
	r := strings.NewReplacer(" ", "+", "/", "~")
	return r.Replace(fen)
}

// ToFen конвертирует UsFEN в FEN (без валидации).
func ToFen(usfen string) string {
	r := strings.NewReplacer("+", " ", "~", "/")
	return r.Replace(usfen)
}
