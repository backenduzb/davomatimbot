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
  // editableOption select maydonidagi tanlangan yozuvning o'zini (masalan sinf
  // nomini) shu sahifadan turib tahrirlash imkonini beradi. Galochka
  // belgilanganda select qulflanadi va matn input orqali nom o'zgartiriladi.
  editableOption?: {
    endpoint: string;
    field?: string; // default: "name"
    toggleLabel?: string;
    toggleLabelKey?: string;
    inputLabel?: string;
    inputLabelKey?: string;
  };
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
