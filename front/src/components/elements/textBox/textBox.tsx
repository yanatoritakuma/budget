"use client";

import TextField from "@mui/material/TextField";
import { FocusEventHandler, memo, useState } from "react";
import IconButton from "@mui/material/IconButton";
import InputAdornment from "@mui/material/InputAdornment";
import Visibility from "@mui/icons-material/Visibility";
import VisibilityOff from "@mui/icons-material/VisibilityOff";
import OutlinedInput from "@mui/material/OutlinedInput";
import InputLabel from "@mui/material/InputLabel";
import FormControl from "@mui/material/FormControl";
import Box from "@mui/material/Box";

type Props = {
  label: string;
  value: string;
  onChange: (event: React.ChangeEvent<HTMLInputElement>) => void;
  className?: string;
  fullWidth?: boolean;
  multiline?: boolean;
  rows?: number;
  password?: boolean;
  onBlur?: FocusEventHandler<HTMLInputElement | HTMLTextAreaElement>;
};

export const TextBox = memo(
  ({
    label,
    value,
    onChange,
    className,
    fullWidth,
    multiline,
    rows,
    password,
    onBlur,
  }: Props) => {
    const [displayPs, setDisplayPs] = useState(false);
    const handleMouseDownPassword = (
      event: React.MouseEvent<HTMLButtonElement>
    ) => {
      event.preventDefault();
    };

    return (
      <>
        {!password ? (
          <TextField
            variant="outlined"
            label={label}
            className={className}
            value={value}
            onChange={onChange}
            fullWidth={fullWidth}
            multiline={multiline}
            rows={rows}
            onBlur={onBlur}
          />
        ) : (
          <FormControl
            variant="outlined"
            sx={{
              maxWidth: "100%",
              width: "100%",
              background: "#fff",
            }}
          >
            <InputLabel>{label}</InputLabel>
            <Box
              sx={{
                maxWidth: "100%",
              }}
            >
              <OutlinedInput
                type={displayPs ? "text" : "password"}
                value={value}
                onChange={onChange}
                label={label}
                endAdornment={
                  <InputAdornment position="end">
                    <IconButton
                      onClick={() => setDisplayPs(!displayPs)}
                      onMouseDown={handleMouseDownPassword}
                      edge="end"
                    >
                      {displayPs ? <VisibilityOff /> : <Visibility />}
                    </IconButton>
                  </InputAdornment>
                }
                fullWidth={fullWidth}
                onBlur={onBlur}
              />
            </Box>
          </FormControl>
        )}
      </>
    );
  }
);

TextBox.displayName = "TextBox";
