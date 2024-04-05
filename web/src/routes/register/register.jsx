import "./register.scss";
import { Link, useNavigate } from "react-router-dom";
import { useState } from "react";
import idmServiceApi from "../../lib/api";
import { IdmCreateAccountRequest } from "../../api/src";

function Register() {
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const navigate = useNavigate();

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError("");
    setIsLoading(true);
    const formData = new FormData(e.target);

    const request = new IdmCreateAccountRequest();
    request.accountName = formData.get("accountName");
    request.password = formData.get("password");

    try {
      idmServiceApi.idmServiceCreateAccount(request);

      navigate("/login");
    } catch (err) {
      setError(err.response.data);
    } finally {
      setIsLoading(false);
    }
  };
  return (
    <div className="registerPage">
      <div className="formContainer">
        <form onSubmit={handleSubmit}>
          <h1>Create an Account</h1>
          <input name="accountName" type="text" placeholder="Account Name" />
          <input name="password" type="password" placeholder="Password" />
          <button disabled={isLoading}>Register</button>
          {error && <span>{error}</span>}
          <Link to="/login">Do you have an account?</Link>
        </form>
      </div>
    </div>
  );
}

export default Register;
