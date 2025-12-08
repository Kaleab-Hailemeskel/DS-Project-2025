'use client'
import { motion, AnimatePresence } from "framer-motion";
import { useState } from "react";
import img from "../../assets/sabrina.png"
import { subtitle, title } from "@/components/primitives";

const SideBySideComponent = () => {
  const [isVisible, setIsVisible] = useState(true);

  // Variants for slide-fade from right (to left direction)
  const variants = {
    hidden: { opacity: 0, x: 50 }, // Start invisible and 50px to the right
    visible: { opacity: 1, x: 0 }, // Fade in and slide left to position
  };

  return (
    <div className="p-4 mx-auto">
      {/* Toggle for demo; remove if not needed, but enables exit animation */}
      <button 
        className="mb-4 px-4 py-2 bg-blue-500 text-white rounded"
        onClick={() => setIsVisible(!isVisible)}
      >
        Toggle Visibility
      </button>
      
      <AnimatePresence>
        {isVisible && (
          <motion.div 
            className="flex flex-row items-center shadow-md rounded-lg overflow-hidden w-full"
            initial="hidden"
            animate="visible"
            exit="hidden"
            variants={variants}
            transition={{ duration: 0.5, ease: "easeOut" }} // Parent container animation (optional, for cohesion)
          >
            {/* Image with right-fade gradient */}
           {/* Image with right-fade gradient */}
<motion.div
  className="relative w-full"  // ensure width is fixed
  initial="hidden"
  animate="visible"
  exit="hidden"
  variants={variants}
  transition={{ duration: 0.4, ease: "easeOut" }}
>
  <img 
    src={img.src}
    alt="Descriptive Image"
    className="w-full h-full object-cover"
  />

  {/* RIGHT-SIDE GRADIENT FADE */}
  <div className="absolute inset-0 pointer-events-none bg-gradient-to-r from-transparent to-black/100" />
</motion.div>

            
            {/* Text div on right */}
            <motion.div
              className="w-1/2 p-4"
              initial="hidden"
              animate="visible"
              exit="hidden"
              variants={variants}
              transition={{ duration: 0.6, ease: "easeOut" }} // Text animates slower (0.6s)
            >
               <div className="inline-block max-w-xl text-center justify-center">
                      <span className={title()}>Explore&nbsp;</span>
                      <span className={title({ color: "pink" })}>Music&nbsp;</span>
                      <br />
                      <span className={title()}>
                        like never before
                      </span>
                      <div className={subtitle({ class: "mt-4" })}>
                        Pop, Rap and everything in between
                      </div>
                    </div>
            </motion.div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
};

export default SideBySideComponent;