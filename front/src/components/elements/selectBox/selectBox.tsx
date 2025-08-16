"use client";

import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import Select, { SelectChangeEvent } from "@mui/material/Select";

type Option = {
  value: string | number;
  label: string;
};

type Props = {
  label: string;
  value: string | number;
  onChange: (event: SelectChangeEvent<string | number>) => void;
  options: Option[];
  name?: string;
  required?: boolean;
};

export const SelectBox = ({
  label,
  value,
  onChange,
  options,
  name,
  required,
}: Props) => {
  const labelId = `${name}-label`;

  return (
    <FormControl fullWidth variant="outlined">
      <InputLabel id={labelId}>{label}</InputLabel>
      <Select
        labelId={labelId}
        id={name}
        name={name}
        value={value}
        label={label}
        onChange={onChange}
        required={required}
      >
        {options.map((option) => (
          <MenuItem key={option.value} value={option.value}>
            {option.label}
          </MenuItem>
        ))}
      </Select>
    </FormControl>
  );
};