// type Photo struct {
// 	Id        uuid.UUID `json:"id" db:"photo_uuid"`
// 	FileName  string    `json:"file_name" db:"photo_filename"`
// 	Width     int       `json:"width" db:"photo_width"`
// 	Height    int       `json:"height" db:"photo_height"`
// 	Timestamp time.Time `json:"timestamp" db:"photo_uploaded"`
// 	Type      PhotoType `json:"type" db:"photo_type"`
// 	Uploader  int       `json:"user_id" db:"photo_uploader"`
// }

export interface Photo {
  id: string;
  file_name: string;
  width: number;
  height: number;
}
