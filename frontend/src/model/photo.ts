export interface Format {
  purpose: number;
  format: string;
  width: number;
  height: number;
  size: number;
}

export interface Photo {
  id: string;
  file_name: string;
  formats: Array<Format>;
}

export const PurposePreview = 1;
export const PurposeScreen = 2;
export const PurposeDownload = 3;
