"use client";

import Button from "@mui/material/Button";
import { memo } from "react";
import IconButton from "@mui/material/IconButton";
import PhotoCamera from "@mui/icons-material/PhotoCamera";

type Props = {
  children?: string;
  onClick?: () => void;
  upload?: boolean;
  onChange?: (e: React.ChangeEvent<HTMLInputElement>) => void;
  variant?: "text" | "contained" | "outlined";
  disabled?: boolean;
  size?: "small" | "medium" | "large";
};

export const ButtonBox = memo((props: Props) => {
  const { children, onClick, upload, onChange, variant, disabled, size } =
    props;

  return (
    <>
      {!upload ? (
        <Button
          onClick={onClick}
          variant={variant !== undefined ? variant : "contained"}
          disabled={disabled}
          size={size}
        >
          {children}
        </Button>
      ) : (
        <IconButton
          color="primary"
          aria-label="upload picture"
          component="label"
        >
          <input hidden accept="image/*" type="file" onChange={onChange} />
          <PhotoCamera />
        </IconButton>
      )}
    </>
  );
});

ButtonBox.displayName = "ButtonBox";
