type FocusIndicatorSettings = {
  type: "bar" | "underline";
  ignoreParent?: boolean;
}

type Validity = {
  valid: boolean;
  message: string;
}

type InputValidation = (value: string) => Promise<Validity>;
type FileValidation = (value: FileList) => Promise<Validity>;

type CacheEntry<T> = {
  date: number;
  value: T | null;
}