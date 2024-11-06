type Option = {
  value: string;
  name: string;
}

type FocusIndicatorSettings = {
  type: "bar" | "underline";
  ignoreParent?: boolean;
}

type Validity = {
  valid: boolean;
  message: string;
}

type InputValidation = (value: string) => Validity;