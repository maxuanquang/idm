import { useContext, useEffect, useState } from "react";
import "./homePage.scss";
import { AuthContext } from "../../context/AuthContext";
import idmServiceApi from "../../lib/api";

function HomePage() {
  const { currentUser } = useContext(AuthContext);

  const [downloadTaskList, setDownloadTaskList] = useState([]);
  const [page, setPage] = useState(1);
  const [limit, setLimit] = useState(10);
  const [totalDownloadTaskCount, setTotalDownloadTaskCount] = useState(0);

  useEffect(() => {
    if (currentUser) {
      getDownloadTaskList();
    }
  }, [page]);

  async function getDownloadTaskList() {
    try {
      idmServiceApi.idmServiceGetDownloadTaskList(
        {
          offset: (page - 1) * limit,
          limit: limit,
        },
        (err, data, response) => {
          if (err != null) {
            console.log(err);
          } else {
            setDownloadTaskList([...data.downloadTaskList]);
            setTotalDownloadTaskCount(data.totalDownloadTaskCount);
          }
        }
      );
    } catch (err) {
      console.log(err);
    }
  }

  function loadNextPage() {
    setPage(page + 1);
  }

  function loadPreviousPage() {
    setPage(page - 1);
  }

  async function handleDownloadFile(task) {
    
    try {
      const request = idmServiceApi.idmServiceGetDownloadTaskFile(
        task.id,
        (err, data, response) => {
          console.log(err);
          console.log(atob(data.result.data));
          console.log(response);

          const blob = new Blob([atob(data.result.data)]); // Create a Blob object

          // Create an <a> element to trigger the download
          const link = document.createElement("a");
          link.href = window.URL.createObjectURL(blob);
          link.download = "download";
          link.click();
        }
      );
      console.log(request);
    } catch (err) {
      console.error("Error downloading file:", err);
    }
  }

  const currentPage = page;
  const totalPages = Math.ceil(totalDownloadTaskCount / limit);

  return (
    <div className="homePage">
      {currentUser ? (
        <>
          {/* Render download task list here */}
          <ul>
            {downloadTaskList.map((task) => (
              <li key={task.id}>
                <div>ID: {task.id}</div>
                <div>URL: {task.url}</div>
                <div>Status: {task.downloadStatus}</div>
                {/* Render other properties of IdmDownloadTask */}
                <button onClick={() => handleDownloadFile(task)}>
                  Get Downloaded File
                </button>
              </li>
            ))}
          </ul>
          {/* Previous Page button */}
          {page > 1 && (
            <button onClick={loadPreviousPage}>Load Previous Page</button>
          )}
          {/* Load More button */}
          {page < totalPages && (
            <button onClick={loadNextPage}>Load Next Page</button>
          )}
          {/* Display current page and total pages */}
          <div>
            Current Page: {currentPage} / Total Pages: {totalPages}
          </div>
        </>
      ) : (
        <></>
      )}
    </div>
  );
}

export default HomePage;
