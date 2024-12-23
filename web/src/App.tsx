import { Link } from "@mui/material";
import Button from "@mui/material/Button/Button";
import TextField from "@mui/material/TextField/TextField";
import { useEffect, useState } from "react";
import { BrowserRouter, Route, Routes, useLocation } from "react-router";
import "./App.css";
import { CREATE_LINK, GET_LINK_BY_HASH } from "./endpoints";
import type { ResponseCreateLink } from "./type";

function CreateShortLinkPage() {
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
              href={`http://158.255.74.123/${result}`}
              target="_blank"
            >{`http://158.255.74.123/${result}`}</Link>
          </div>
        )}
      </div>
    </div>
  );
}

function HashToOriginalUrlPage() {
  const location = useLocation();

  const pathname = location.pathname.startsWith("/")
    ? location.pathname.slice(1)
    : location.pathname;

  useEffect(() => {
    fetch(GET_LINK_BY_HASH + `?hash=${pathname}`, {
      method: "get",
      headers: {
        Accept: "application/json",
        "Content-Type": "application/json",
      },
    })
      .then((res) => res.json())
      .then((res: { URL: string }) => {
        window.location.href = res.URL;
      })
      .catch((err) => {
        window.console.log("err = ", err);
      });
  }, [pathname]);
  return null;
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/:hash" element={<HashToOriginalUrlPage />} />
        <Route index path="*" element={<CreateShortLinkPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
