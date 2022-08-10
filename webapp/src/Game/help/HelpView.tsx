import React, { useEffect } from "react";
import { useMedia } from "react-use";
import { Navigate, Route, Routes, useLocation } from "react-router-dom";
import classnames from "classnames";
import { useStore } from "../../store";
import { HelpMenu } from "./HelpMenu";
import { GettingStarted } from "./GettingStarted";
import { Buildings } from "./buildings/Buildings";
import { Kitchen } from "./buildings/Kitchen";
import { Chef } from "./educations/Chef";
import { Shop } from "./buildings/Shop";
import { TownCentre } from "./buildings/TownCentre";
import { House } from "./buildings/House";
import { School } from "./buildings/School";
import { MarketingHq } from "./buildings/MarketingHQ";
import { ResearchInstitute } from "./buildings/ResearchInstitute";
import { Educations } from "./educations/Educations";
import { Salesmouse } from "./educations/Salesmouse";
import { Guard } from "./educations/Guard";
import { Thief } from "./educations/Thief";
import { Publicist } from "./educations/Publicist";

const HELP_QUEST_ID = "6";

const useIsMinXl = () => useMedia("(min-width: 1280px)", false);

function useScrollHashElementIntoView() {
  const { pathname, hash, key } = useLocation();
  useEffect(() => {
    if (hash === "") {
      window.scrollTo(0, 0);
    } else {
      setTimeout(() => {
        const id = hash.replace("#", "");
        const element = document.getElementById(id);
        if (element) {
          element.scrollIntoView();
        }
      }, 0);
    }
  }, [pathname, hash, key]);
}

function useCompleteHelpQuest() {
  const completeQuest = useStore((state) => state.completeQuest);
  const helpQuestState = useStore(
    (state) => state.gameState.quests[HELP_QUEST_ID]
  );

  useEffect(() => {
    if (helpQuestState && !helpQuestState.completed) {
      completeQuest(HELP_QUEST_ID);
    }
  }, [helpQuestState, completeQuest]);
}

const HelpView: React.VFC<{}> = () => {
  useCompleteHelpQuest();
  useScrollHashElementIntoView();

  const isMinXl = useIsMinXl();

  return (
    <div
      className={classnames("flex", "items-center", "flex-col", "mt-2", "p-2")}
      data-cy="game-help"
    >
      <h2>Game Help</h2>

      <div
        className={classnames("flex gap-4", {
          "flex-col": !isMinXl,
          "items-baseline": isMinXl,
          "items-center": !isMinXl,
        })}
      >
        <HelpMenu className="shrink-0 sticky top-4" alwaysExpanded={isMinXl} />
        <div className="grow-1 min-w-[65ch]">
          <Routes>
            <Route index element={<Navigate to="getting-started" replace />} />
            <Route path="getting-started" element={<GettingStarted />} />
            <Route path="buildings" element={<Buildings />} />
            <Route path="buildings/town-centre" element={<TownCentre />} />
            <Route path="buildings/kitchen" element={<Kitchen />} />
            <Route path="buildings/shop" element={<Shop />} />
            <Route path="buildings/house" element={<House />} />
            <Route path="buildings/school" element={<School />} />
            <Route path="buildings/marketing-hq" element={<MarketingHq />} />
            <Route
              path="buildings/research-institute"
              element={<ResearchInstitute />}
            />
            <Route path="educations" element={<Educations />} />
            <Route path="educations/chef" element={<Chef />} />
            <Route path="educations/salesmouse" element={<Salesmouse />} />
            <Route path="educations/guard" element={<Guard />} />
            <Route path="educations/thief" element={<Thief />} />
            <Route path="educations/publicist" element={<Publicist />} />
          </Routes>
        </div>
      </div>
    </div>
  );
};

export default HelpView;
