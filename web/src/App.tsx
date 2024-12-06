import Button from "@mui/material/Button/Button";
import TextField from "@mui/material/TextField/TextField";
import { useState } from "react";
import "./App.css";
import { CREATE_LINK } from "./endpoints";
import type { ResponseCreateLink } from "./type";
import { Link } from "@mui/material";

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
      .then((res) => res.json())
      .then((response: ResponseCreateLink) => setResult(response.hash))
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
        {result && (
          <div className="Result">
            <Link
              href={`http://localhost/${result}`}
            >{`http://localhost/${result}`}</Link>
          </div>
        )}
      </div>
    </div>
  );
}

export default App;
