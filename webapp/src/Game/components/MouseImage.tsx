import React from "react";
import classnames from "classnames";

import { isNonNullable } from "../../utils";

import { ReactComponent as SvgArmsBasic } from "images/parts/arms_Basic.svg";
import { ReactComponent as SvgBodiesMouse1 } from "images/parts/bodies_Mouse1.svg";
import { ReactComponent as SvgEyesEyes1 } from "images/parts/eyes_Eyes1.svg";
import { ReactComponent as SvgEyesEyes2 } from "images/parts/eyes_Eyes2.svg";
import { ReactComponent as SvgEyesEyes3 } from "images/parts/eyes_Eyes3.svg";
import { ReactComponent as SvgFeetBasic } from "images/parts/feet_Basic.svg";
import { ReactComponent as SvgFeetBigFeet1 } from "images/parts/feet_BigFeet1.svg";
import { ReactComponent as SvgFeetSmallFeet1 } from "images/parts/feet_SmallFeet1.svg";
import { ReactComponent as SvgMouthsMixed1 } from "images/parts/mouths_Mixed1.svg";
import { ReactComponent as SvgMouthsSmile1 } from "images/parts/mouths_Smile1.svg";
import { ReactComponent as SvgMouthsSmile2 } from "images/parts/mouths_Smile2.svg";
import { ReactComponent as SvgMouthsSmile3 } from "images/parts/mouths_Smile3.svg";
import { ReactComponent as SvgTailsTail1 } from "images/parts/tails_Tail1.svg";
import { ReactComponent as SvgTailsTail2 } from "images/parts/tails_Tail2.svg";
import { ReactComponent as SvgTailsTail3 } from "images/parts/tails_Tail3.svg";
import { ReactComponent as SvgTailsTail4 } from "images/parts/tails_Tail4.svg";
import { ReactComponent as SvgTailsTail5 } from "images/parts/tails_Tail5.svg";
import { ReactComponent as SvgEyesExtra1EyeCover1 } from "images/parts/eyes_extra1_EyeCover1.svg";
import { ReactComponent as SvgEyesExtra1EyeStars1 } from "images/parts/eyes_extra1_EyeStars1.svg";
import { ReactComponent as SvgEyesExtra2Glasses1 } from "images/parts/eyes_extra2_Glasses1.svg";
import { ReactComponent as SvgEyesExtra2Glasses2 } from "images/parts/eyes_extra2_Glasses2.svg";
import { ReactComponent as SvgEyesExtra2EyePatch1 } from "images/parts/eyes_extra2_EyePatch1.svg";
import { ReactComponent as SvgOutfitsThiefOutfit1 } from "images/parts/outfits_ThiefOutfit1.svg";
import { ReactComponent as SvgOutfitsOutfit1 } from "images/parts/outfits_Outfit1.svg";
import { ReactComponent as SvgOutfitsGuardOutfit1 } from "images/parts/outfits_GuardOutfit1.svg";
import { ReactComponent as SvgHatsRedHat1 } from "images/parts/hats_RedHat1.svg";
import { ReactComponent as SvgHatsThiefHat1 } from "images/parts/hats_ThiefHat1.svg";
import { ReactComponent as SvgHatsChefHat1 } from "images/parts/hats_ChefHat1.svg";
import { ReactComponent as SvgHatsGuardHat1 } from "images/parts/hats_GuardHat1.svg";
import { ReactComponent as SvgHatsHat1 } from "images/parts/hats_Hat1.svg";
import { ReactComponent as SvgHatsBucketHat1 } from "images/parts/hats_BucketHat1.svg";
import { ReactComponent as SvgHatsCap1 } from "images/parts/hats_Cap1.svg";
import {
  AppearanceCategory,
  MouseAppearance,
  MouseAppearance_PartRef,
} from "../../generated/appearance";

type PartInfo = {
  component: React.ComponentType<{
    x?: string | number;
    y?: string | number;
    className?: string;
  }>;
  x: number;
  y: number;
};

const partInfosById: Record<string, PartInfo | undefined> = {
  redHat1: {
    component: SvgHatsRedHat1,
    x: 25,
    y: 15,
  },
  thiefHat1: {
    component: SvgHatsThiefHat1,
    x: 25,
    y: 20,
  },
  chefHat1: {
    component: SvgHatsChefHat1,
    x: 16,
    y: 4,
  },
  guardHat1: {
    component: SvgHatsGuardHat1,
    x: 22,
    y: 14,
  },
  hat1: {
    component: SvgHatsHat1,
    x: 19.5,
    y: 13,
  },
  bucketHat1: {
    component: SvgHatsBucketHat1,
    x: 15.5,
    y: 20,
  },
  cap1: {
    component: SvgHatsCap1,
    x: 21.5,
    y: 19.5,
  },
  basic: {
    component: SvgArmsBasic,
    x: 30,
    y: 112,
  },
  basicFeet1: {
    component: SvgFeetBasic,
    x: 4,
    y: 176,
  },
  bigFeet1: {
    component: SvgFeetBigFeet1,
    x: 4,
    y: 176,
  },
  smallFeet1: {
    component: SvgFeetSmallFeet1,
    x: 4,
    y: 176,
  },
  mixedSmile1: {
    component: SvgMouthsMixed1,
    x: 32,
    y: 70,
  },
  smile1: {
    component: SvgMouthsSmile1,
    x: 35,
    y: 70,
  },
  smile2: {
    component: SvgMouthsSmile2,
    x: 32,
    y: 70,
  },
  smile3: {
    component: SvgMouthsSmile3,
    x: 30,
    y: 70,
  },
  tail1: {
    component: SvgTailsTail1,
    x: 35,
    y: 175,
  },
  tail2: {
    component: SvgTailsTail2,
    x: 35,
    y: 175,
  },
  tail3: {
    component: SvgTailsTail3,
    x: 35,
    y: 160,
  },
  tail4: {
    component: SvgTailsTail4,
    x: 35,
    y: 153,
  },
  tail5: {
    component: SvgTailsTail5,
    x: 35,
    y: 175,
  },
  glasses2: {
    component: SvgEyesExtra2Glasses2,
    x: 12,
    y: 50,
  },
  glasses1: {
    component: SvgEyesExtra2Glasses1,
    x: 12,
    y: 50,
  },
  eyePatch1: {
    component: SvgEyesExtra2EyePatch1,
    x: 9.85,
    y: 45.5,
  },
  eyeCover1: {
    component: SvgEyesExtra1EyeCover1,
    x: 10,
    y: 47,
  },
  eyeStars1: {
    component: SvgEyesExtra1EyeStars1,
    x: 17,
    y: 45,
  },
  eyes1: {
    component: SvgEyesEyes1,
    x: 23,
    y: 50,
  },
  eyes2: {
    component: SvgEyesEyes2,
    x: 24,
    y: 51,
  },
  eyes3: {
    component: SvgEyesEyes3,
    x: 27.8,
    y: 52,
  },
  thiefOutfit1: {
    component: SvgOutfitsThiefOutfit1,
    x: 5,
    y: 87,
  },
  guardOutfit1: {
    component: SvgOutfitsGuardOutfit1,
    x: 6.9,
    y: 80.5,
  },
  outfit1: {
    component: SvgOutfitsOutfit1,
    x: 6.9,
    y: 80.5,
  },
  mouse1: {
    component: SvgBodiesMouse1,
    x: 0,
    y: 20,
  },
};

export function MouseImagePart({
  id,
  className,
  withOffset = true,
}: {
  id: string;
  className?: string;
  withOffset?: boolean;
}) {
  const part = partInfosById[id];
  if (!part) {
    console.warn(`Failed to find part with id: ${id}`);
    return null;
  }

  const Component = part.component;

  return (
    <Component
      x={withOffset ? part.x : undefined}
      y={withOffset ? part.y : undefined}
      className={className}
    />
  );
}

const appearanceOrder = [
  AppearanceCategory.TAIL,
  AppearanceCategory.FEET,
  AppearanceCategory.BODY,
  AppearanceCategory.OUTFIT,
  AppearanceCategory.MOUTH,
  AppearanceCategory.EYES_EXTRA1,
  AppearanceCategory.EYES,
  AppearanceCategory.EYES_EXTRA2,
  AppearanceCategory.HAT,
  AppearanceCategory.LEFT_HAND,
  AppearanceCategory.RIGHT_HAND,
];

function getDefaultPart(
  category: AppearanceCategory
): MouseAppearance_PartRef | null {
  switch (category) {
    case AppearanceCategory.BODY:
      return MouseAppearance_PartRef.create({ id: "mouse1" });
    case AppearanceCategory.FEET:
      return MouseAppearance_PartRef.create({ id: "basicFeet1" });
    case AppearanceCategory.MOUTH:
      return MouseAppearance_PartRef.create({ id: "smile1" });
    case AppearanceCategory.EYES:
      return MouseAppearance_PartRef.create({ id: "eyes1" });
    case AppearanceCategory.LEFT_HAND:
      return MouseAppearance_PartRef.create({ id: "basic" });
    default:
      return null;
  }
}

export function MouseImage({
  appearance = MouseAppearance.create(),
  shiftRight,
  className,
  ...rest
}: {
  appearance?: MouseAppearance;
  shiftRight?: boolean;
} & React.SVGProps<SVGSVGElement>) {
  return (
    <svg
      width={320}
      height={400}
      viewBox="0 0 160 200"
      className={classnames(
        {
          "translate-x-4 xs:translate-x-16": shiftRight,
        },
        className
      )}
      {...rest}
    >
      {appearanceOrder
        .map(
          (category) => appearance.parts[category] ?? getDefaultPart(category)
        )
        .filter(isNonNullable)
        .map((part) => (
          <MouseImagePart key={part.id} id={part.id} className="absolute" />
        ))}
    </svg>
  );
}
