import "./login.scss";
import { Link, useNavigate } from "react-router-dom";
import { useState, useContext } from "react";
import { AuthContext } from "../../context/AuthContext";
import idmServiceApi from "../../lib/api";
import { IdmCreateSessionRequest } from "../../api/src";

function Login() {
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const { updateUser } = useContext(AuthContext);

  const navigate = useNavigate();

  const handleFormSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    setError("");

    const formData = new FormData(e.target);

    const accountName = formData.get("accountName");
    const password = formData.get("password");

    const request = new IdmCreateSessionRequest();
    request.accountName = accountName;
    request.password = password;

    try {
      idmServiceApi.idmServiceCreateSession(request, (err, data, response) => {
        if (err != null) {
          setError("wrong username or password");
        } else {
          updateUser(data.account);

          navigate("/");
        }
      });
    } catch (err) {
      setError(err.response.data);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="login">
      <div className="formContainer">
        <form onSubmit={handleFormSubmit}>
          <h1>Welcome back</h1>
          <input
            name="accountName"
            required
            minLength={3}
            maxLength={20}
            type="text"
            placeholder="Account Name"
          />
          <input
            name="password"
            type="password"
            required
            placeholder="Password"
          />
          <button disabled={isLoading}>Login</button>
          {error && <span>{error}</span>}
          <Link to="/register">{"Don't"} you have an account?</Link>
        </form>
      </div>
    </div>
  );
}

export default Login;
