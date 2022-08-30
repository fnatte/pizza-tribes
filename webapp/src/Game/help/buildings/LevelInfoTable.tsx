import React from "react";
import classnames from "classnames";
import { useStore } from "../../../store";
import {
  CellContext,
  createColumnHelper,
  flexRender,
  getCoreRowModel,
  useReactTable,
} from "@tanstack/react-table";
import { Building, BuildingInfo_LevelInfo } from "../../../generated/building";
import { formatNumber } from "../../../utils";

type Row = BuildingInfo_LevelInfo & { level: number };

const columnHelper = createColumnHelper<Row>();

const getFormattedNumber = (
  info: CellContext<Row, number | undefined> | CellContext<Row, number>
): string | null => {
  const value = info.getValue();
  if (value === undefined) {
    return null;
  }
  return formatNumber(value);
};

const columns = [
  columnHelper.accessor("level", {
    id: "level",
    cell: (info) => <span className="font-bold">{info.getValue()}</span>,
    header: "Level",
  }),
  columnHelper.accessor((row) => row.employer?.maxWorkforce, {
    id: "workforce",
    cell: getFormattedNumber,
    header: "Workforce",
  }),
  columnHelper.accessor((row) => row.residence?.beds, {
    id: "beds",
    cell: getFormattedNumber,
    header: "Beds",
  }),
  columnHelper.accessor("cost", {
    cell: getFormattedNumber,
    header: "Cost",
  }),
  columnHelper.accessor("constructionTime", {
    cell: getFormattedNumber,
    header: "Construction Time",
  }),
];

export function LevelInfoTable({
  building,
  className,
}: {
  building: Building;
  className?: string;
}) {
  const gameData = useStore((state) => state.gameData);

  const data =
    gameData?.buildings[building]?.levelInfos
      .map((levelInfo, i) => ({ ...levelInfo, level: i + 1 }))
      .filter((row) => row.cost > 0 || row.constructionTime > 0) ?? [];

  const table = useReactTable({
    data,
    columns,
    getCoreRowModel: getCoreRowModel(),
    state: {
      columnVisibility: {
        workforce: data.some((row) => row.employer?.maxWorkforce),
        beds: data.some((row) => row.residence?.beds),
      },
    },
  });

  return (
    <div className={classnames("bg-green-50 p-2 xs:p-4", className)}>
      <table className="m-0">
        <thead>
          {table.getHeaderGroups().map((headerGroup) => (
            <tr key={headerGroup.id}>
              {headerGroup.headers.map((header) => (
                <th key={header.id}>
                  {header.isPlaceholder
                    ? null
                    : flexRender(
                        header.column.columnDef.header,
                        header.getContext()
                      )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {table.getRowModel().rows.map((row) => (
            <tr key={row.id}>
              {row.getVisibleCells().map((cell) => (
                <td key={cell.id}>
                  {flexRender(cell.column.columnDef.cell, cell.getContext())}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
