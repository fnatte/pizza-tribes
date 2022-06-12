import React from "react";
import classnames from "classnames";
import { ReactComponent as SvgTownCentre } from "../../../images/town-centre.svg";
import { useStore } from "../../store";
import PopulationTable from "../PopulationTable";
import { Link } from "react-router-dom";
import { smallPrimaryButton } from "../../styles";

function MouseTable({ className }: { className?: string }) {
  const mice = useStore((state) => state.gameState.mice);
  const gameData = useStore((state) => state.gameData);
  const educations = gameData?.educations ?? {};

  return (
    <table className={className}>
      {
        <thead>
          <tr>
            <th colSpan={2} className="p-2 text-left text-sm text-gray-700">
              Mice
            </th>
          </tr>
        </thead>
      }
      <tbody>
        {Object.values(mice).length === 0 && (
          <tr>
            <td colSpan={2} className="p-2 w-80 text-left text-gray-700">
              When your town has some mice, they will show up here.
            </td>
          </tr>
        )}
        {Object.entries(mice).map(([id, mouse]) => {
          return (
            <tr key={id}>
              <td className="p-2">{mouse.name}</td>
              <td className="p-2">
                {!mouse.isEducated
                  ? "Uneducated"
                  : mouse.isBeingEducated
                  ? "In school"
                  : educations[mouse.education].title}
              </td>
              <td className="p-2">
                <Link to={`/mouse/${id}`}>
                  <button className={classnames(smallPrimaryButton)}>
                    Visit
                  </button>
                </Link>
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

function PopulationSection() {
  return (
    <section className="my-6">
      <h3>Population</h3>

      <div className="prose my-6 text-gray-700">
        <p>Your tribesmice have good morale and is working efficiently.</p>
      </div>

      <div className="flex flex-wrap items-start gap-24">
        <PopulationTable showHeader showZeroes showTotalCount />
        <MouseTable />
      </div>
    </section>
  );
}

export function TownCentre() {
  return (
    <div className={classnames("container", "px-2", "max-w-2xl", "mb-8")}>
      <h2>Town Centre</h2>
      <div className="flex gap-6 items-center">
        <SvgTownCentre height={100} width={100} />
        <div className="prose my-6 text-gray-700">
          <p>
            This is where tribesmice go to decide on important tribe issues.
          </p>
        </div>
      </div>
      <PopulationSection />
    </div>
  );
}
