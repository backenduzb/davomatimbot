export type Column = {
  key: string;
  label?: string;
  labelKey?: string;
  formatter?: (value: unknown, row: Record<string, unknown>) => string;
  link?: (row: Record<string, unknown>) => string;
};

export type Field = {
  key: string;
  label?: string;
  labelKey?: string;
  type:
    | "text"
    | "number"
    | "date"
    | "email"
    | "password"
    | "select"
    | "checkbox"
    | "file";
  required?: boolean;
  optionsEndpoint?: string;
  optionLabel?: string;
  optionValue?: string;
  options?: {
    label?: string;
    labelKey?: string;
    value: string | number;
  }[];
};

export type Filter = {
  key: string;
  label?: string;
  labelKey?: string;
  type: "select" | "date";
  options?: {
    label?: string;
    labelKey?: string;
    value: string | number;
  }[];
};
