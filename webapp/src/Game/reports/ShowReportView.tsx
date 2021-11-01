import React, { useCallback, useEffect } from "react";
import { classnames } from "tailwindcss-classnames";
import { useParams } from "react-router-dom";
import { useStore } from "../../store";
import { formatISO9075 } from "date-fns";
import { parseDateNano } from "../../utils";
import ReactMarkdown from "react-markdown";

function ShowReportView() {
  const { id } = useParams();

  const readReport = useStore((state) => state.readReport);

  const report = useStore(
    useCallback(
      (state) => state.reports.filter((report) => report.id === id)[0],
      [id]
    )
  );

  useEffect(() => {
    if (report?.id) {
      readReport(report.id);
    }
  }, [report?.id]);

  return (
    <div
      className={classnames(
        "flex",
        "flex-col",
        "items-center",
        "justify-center",
        "mt-2",
        "p-2"
      )}
    >
      <h2>Report</h2>
      <div
        className={classnames(
          "container",
          "mx-auto",
          "mt-4",
          "p-4",
          "max-w-md",
          "bg-green-50"
        )}
      >
        {report && (
          <>
            <h3>{report.title}</h3>
            <div>{formatISO9075(parseDateNano(report.createdAt))}</div>
            <ReactMarkdown className={classnames("prose" as any, "mt-6", "text-gray-700")} disallowedElements={['img']}>
              {report.content}
            </ReactMarkdown>
          </>
        )}
      </div>
    </div>
  );
}

export default ShowReportView;
