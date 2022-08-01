import React, { useRef, useState } from "react";
import { produce } from "immer";
import classnames from "classnames";
import { ReactComponent as SvgArrowRight } from "images/icons/arrow-right.svg";
import { button } from "../../styles";
import {
  AppearanceCategory,
  MouseAppearance,
} from "../../generated/appearance";
import { useStore } from "../../store";
import { MouseImage, MouseImagePart } from "../components/MouseImage";

function GallerySection({
  title,
  category,
  expanded,
  onExpandToggle,
  className,
  selected,
  onSelect,
  allowEmpty,
}: {
  title: string;
  category: AppearanceCategory;
  expanded: boolean;
  onExpandToggle?: () => void;
  className?: string;
  selected?: string | null;
  onSelect?: (value: string | null) => void;
  allowEmpty?: boolean;
}) {
  const appearanceParts = useStore((state) => state.gameData?.appearanceParts);
  const sectionParts = Object.values(appearanceParts ?? {}).filter(
    (part) => part.category === category
  );

  return (
    <section className={className}>
      <button
        type="button"
        onClick={onExpandToggle}
        className="w-full bg-green-200"
        data-cy="gallery-section-title"
        aria-expanded={expanded}
      >
        {title}
      </button>
      {expanded && (
        <div className="flex gap-8 flex-wrap mt-2 p-1">
          {allowEmpty && (
            <button
              type="button"
              className={classnames("bg-green-50 p-4 border", {
                outline: !selected,
              })}
              onClick={() => onSelect?.(null)}
              data-cy="gallery-section-item-nothing"
            >
              Nothing
            </button>
          )}
          {sectionParts.map((part) => {
            return (
              <button
                type="button"
                key={part.id}
                className={classnames("bg-green-50 p-4", {
                  outline: selected === part.id,
                })}
                onClick={() => onSelect?.(part.id)}
                data-cy="gallery-section-item"
              >
                <MouseImagePart id={part.id} withOffset={false} />
              </button>
            );
          })}
        </div>
      )}
    </section>
  );
}

function PartGallery({
  className,
  appearance,
  onChange,
}: {
  className?: string;
  appearance: MouseAppearance;
  onChange: (value: MouseAppearance) => void;
}) {
  const [currentSection, setCurrectSection] = useState(AppearanceCategory.BODY);

  const sections = [
    { category: AppearanceCategory.BODY, title: "Bodies" },
    { category: AppearanceCategory.OUTFIT, title: "Outfits", allowEmpty: true },
    { category: AppearanceCategory.HAT, title: "Hats", allowEmpty: true },
    { category: AppearanceCategory.EYES, title: "Eyes" },
    {
      category: AppearanceCategory.EYES_EXTRA1,
      title: "Eyes Extra 1",
      allowEmpty: true,
    },
    {
      category: AppearanceCategory.EYES_EXTRA2,
      title: "Eyes Extra 2",
      allowEmpty: true,
    },
    { category: AppearanceCategory.MOUTH, title: "Mouths" },
    { category: AppearanceCategory.FEET, title: "Feet" },
    // { category: AppearanceCategory.ARM, title: "Arms" },
    { category: AppearanceCategory.TAIL, title: "Tails" },
  ] as { category: AppearanceCategory; title: string; allowEmpty?: boolean }[];

  return (
    <div className={className}>
      {sections.map(({ category, title, allowEmpty }) => (
        <GallerySection
          key={category}
          category={category}
          title={title}
          expanded={currentSection === category}
          onExpandToggle={() => setCurrectSection(category)}
          className="mb-4"
          selected={appearance.parts[category]?.id}
          onSelect={(value) =>
            onChange(
              produce(appearance, (draft) => {
                if (value) {
                  draft.parts[category] = {
                    id: value,
                    color: 0,
                  };
                } else {
                  delete draft.parts[category];
                }
              })
            )
          }
          allowEmpty={allowEmpty}
        />
      ))}
    </div>
  );
}

export function MouseEditor({
  className,
  onCancel,
  onSave,
  defaultAppearance,
}: {
  className?: string;
  onCancel?: () => void;
  onSave?: (appearance: MouseAppearance) => void;
  defaultAppearance?: MouseAppearance;
}) {
  const [appearance, setAppearance] = useState<MouseAppearance>(
    () => defaultAppearance ?? MouseAppearance.create()
  );

  const ref = useRef<HTMLFormElement>(null);
  const [showSlideRight, setShowSlideRight] = useState(true);

  return (
    <form
      className={classnames(
        "flex sm:justify-center gap-8 overflow-x-scroll scroll-smooth snap-x snap-mandatory",
        className
      )}
      onSubmit={(e) => {
        e.preventDefault();
        onSave?.(appearance);
      }}
      ref={ref}
    >
      <div className="min-w-[100vw] sm:min-w-fit shrink-0 sm:shrink-1 snap-always snap-start">
        <div className="relative">
          <MouseImage
            className="w-full mb-2"
            appearance={appearance}
            shiftRight
            onClick={() => setShowSlideRight(false)}
          />
          {showSlideRight && (
            <div className="flex justify-center items-center gap-1 absolute left-1/2 top-1/2 w-3/4 bg-green-50 -translate-x-1/2 -translate-y-1/2 rounded-xl animate-pulse sm:hidden">
              Slide right to change
              <SvgArrowRight className={"h-10 w-10"} />
            </div>
          )}
        </div>
        <div className="mt-1 sm:mt-4 flex justify-center gap-4 ">
          <button
            type="button"
            className={classnames(...button, "bg-gray-600")}
            onClick={onCancel}
            data-cy="mouse-editor-cancel-button"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={false}
            className={classnames(...button, "bg-green-600")}
            data-cy="mouse-editor-save-button"
          >
            Save
          </button>
        </div>
      </div>
      <PartGallery
        className="grow min-w-[100vw] sm:min-w-[auto] sm:max-w-[500px] snap-always snap-end overflow-y-scroll"
        appearance={appearance}
        onChange={(value) => {
          setAppearance(value);
          setShowSlideRight(false);
          ref.current?.scrollTo({ left: 0 });
        }}
      />
    </form>
  );
}
