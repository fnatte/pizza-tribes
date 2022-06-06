import { formatISO9075 } from "date-fns";
import React from "react";
import { Link } from "react-router-dom";
import classnames from "classnames";
import { useStore } from "../../store";
import { parseDateNano } from "../../utils";

function ListReportsView() {
  const reports = useStore((state) => state.reports);

  return (
    <div
      className={classnames("flex", "items-center", "flex-col", "mt-2", "p-2")}
    >
      <table
        className={classnames(
          "w-full",
          "max-w-md",
          "my-4",
          "border-collapse",
          "border-green-400",
          "border-2"
        )}
      >
        <thead>
          <tr>
            <th className={classnames("p-1", "text-left")}>Date</th>
            <th className={classnames("p-1", "pl-8", "text-left")}>Title</th>
          </tr>
        </thead>
        <tbody>
          {reports.map((report, i) => (
            <tr
              key={report.id}
              className={classnames({
                "bg-green-200": i % 2 === 0,
              })}
            >
              <td className={classnames("p-1", "w-1", "whitespace-nowrap")}>
                {formatISO9075(parseDateNano(report.createdAt))}
              </td>
              <td
                className={classnames({
                  "p-1": true,
                  "pl-8": true,
                  "font-bold": report.unread,
                  underline: true,
                })}
              >
                <Link to={`/reports/${report.id}`}>{report.title}</Link>
              </td>
            </tr>
          ))}
          {reports.length === 0 && (
            <tr className={classnames("bg-green-200")}>
              <td colSpan={2} className={classnames("p-1")}>
                You do not have any reports.
              </td>
            </tr>
          )}
        </tbody>
      </table>
    </div>
  );
}

export default ListReportsView;
