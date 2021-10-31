import { newLangController, newSearchController, newToCController } from './controllers';
import Alpine from 'jslibs/alpinejs/v3/alpinejs/dist/module.esm.js';

// Register AlpineJS data controllers.
Alpine.data('searchController', newSearchController);
Alpine.data('tocController', newToCController);
Alpine.data('langController', newLangController);

// Start AlpineJS.
Alpine.start();
