import TextField from "@mui/material/TextField/TextField";
import Button from "@mui/material/Button/Button";
import { useState } from "react";
import "./App.css";
import { CREATE_LINK } from "./endpoints";

function App() {
  const [value, setValue] = useState<string>("");
  const [result, setResult] = useState<string>("");

  const handleOnChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const newValue = e.target.value;
    setValue(newValue);
  };

  const handleSubmit = () => {
    fetch(CREATE_LINK, {
      method: "post",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        URL: value,
      }),
    })
      .then(() => {
        setResult("Link successfuly created");
      })
      .catch((error) => {
        window.console.log("error = ", error);
      });
  };

  return (
    <div className="Main">
      <div className="Section">
        <div className="Form">
          <TextField
            label="Enter your URL"
            value={value}
            onChange={handleOnChange}
          />
          <Button variant="contained" onClick={handleSubmit}>
            Shorten URL
          </Button>
        </div>
        {result && <p>{result}</p>}
      </div>
    </div>
  );
}

export default App;
