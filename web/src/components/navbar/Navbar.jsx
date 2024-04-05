// import { useContext, useState } from "react";
// import "./navbar.scss";
// import { Link, useNavigate } from "react-router-dom";
// import { AuthContext } from "../../context/AuthContext";

// function Navbar() {
//   const [open, setOpen] = useState(false);

//   const { currentUser, updateUser } = useContext(AuthContext);

//   const navigate = useNavigate();

//   function handleLogoutClick() {
//     localStorage.removeItem("user");
//     updateUser(null);
//     navigate("/login");
//   }

//   function handleNewDownloadTaskClick() {
//   }

//   return (
//     <nav>
//       <div className="left">
//         <a href="/" className="logo">
//           <img src="/logo.png" alt="" />
//           <span>Internet Download Manager</span>
//         </a>
//       </div>
//       <div className="right">
//         {currentUser ? (
//           <div className="user">
//             <img src={currentUser.avatar || "/noavatar.jpg"} alt="" />
//             <span>{currentUser.accountName}</span>
//             <button className="logoutButton" onClick={handleLogoutClick}>
//               Logout
//             </button>
//           </div>
//         ) : (
//           <>
//             <a href="/login">Sign in</a>
//             <a href="/register" className="register">
//               Sign up
//             </a>
//           </>
//         )}
//         {currentUser && (
//           <div className="user">
//             <button className="newDownloadTaskButton" onClick={handleNewDownloadTaskClick}>
//               New Download Task
//             </button>
//           </div>
//         )}
//       </div>
//     </nav>
//   );
// }

// export default Navbar;
import { useContext, useState } from "react";
import "./navbar.scss";
import { Link, useNavigate } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import idmServiceApi from "../../lib/api";
import {
  IdmCreateAccountRequest,
  IdmCreateDownloadTaskRequest,
} from "../../api/src";

function Navbar() {
  const [open, setOpen] = useState(false);
  const [downloadType, setDownloadType] = useState("");
  const [url, setUrl] = useState("");
  const [message, setMessage] = useState("");
  const [error, setError] = useState(false);

  const { currentUser, updateUser } = useContext(AuthContext);

  const navigate = useNavigate();

  function handleLogoutClick() {
    localStorage.removeItem("user");
    updateUser(null);
    navigate("/login");
  }

  function handleNewDownloadTaskClick() {
    setOpen(!open);
  }

  async function handleSubmit(event) {
    event.preventDefault();
    try {
      const request = new IdmCreateDownloadTaskRequest();
      request.download_type = downloadType;
      request.url = url;

      idmServiceApi.idmServiceCreateDownloadTask(
        request,
        (err, data, response) => {
          if (err != null) {
            setMessage("Failed to create download task.");
            setError(true);
          } else {
            setMessage("Download task created successfully.");
            setError(false);
          }
        }
      );
    } catch (err) {
      setMessage("Failed to create download task.");
      setError(true);
    }
  }

  return (
    <nav>
      <div className="left">
        <a href="/">
          <span className="logo">Internet Download Manager</span>
        </a>
      </div>
      <div className="right">
        {currentUser ? (
          <div className="user">
            <span className="userName">{currentUser.accountName}</span>
            <button className="logoutButton" onClick={handleLogoutClick}>
              Logout
            </button>
          </div>
        ) : (
          <>
            <a href="/login" className="signinLink">Sign in</a>
            <a href="/register" className="registerLink">
              Sign up
            </a>
          </>
        )}
        {currentUser && (
          <div className="user">
            <button
              className="newDownloadTaskButton"
              onClick={handleNewDownloadTaskClick}
            >
              New Download Task
            </button>
          </div>
        )}
      </div>
      {/* Popup for new download task */}
      {open && (
        <div className="popup">
          <form onSubmit={handleSubmit}>
            <label htmlFor="downloadType">Download Type:</label>
            <input
              type="text"
              id="downloadType"
              value={downloadType}
              onChange={(e) => setDownloadType(e.target.value)}
              required
            />
            <label htmlFor="url">URL:</label>
            <input
              type="text"
              id="url"
              value={url}
              onChange={(e) => setUrl(e.target.value)}
              required
            />
            <button type="submit">Create Download Task</button>
          </form>
          {message && (
            <div className={error ? "error" : "success"}>{message}</div>
          )}
        </div>
      )}
    </nav>
  );
}

export default Navbar;
