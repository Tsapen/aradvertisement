<template style="margin: 0; overflow: hidden;">
    <a-scene vr-mode-ui="enabled: false"
        embedded arjs="sourceType: webcam; debugUIEnabled: false;">
        <a-text
            value="Picture"
            look-at="[gps-camera]"
            scale="45 90 90"
            gps-entity-place="latitude: 55.645811; longitude: 37.649352;"
        ></a-text>
        <a-camera gps-camera rotation-reader> </a-camera>
    </a-scene>
</template>
<script type="text/javascript">
    import Vue from Vue
    export default {
        name: "",
        created() {
            let aframe = document.createElement('script');    
            aframe.setAttribute('src',"https://aframe.io/releases/1.0.4/aframe.min.js");
            document.head.appendChild(aframe);

            let component = document.createElement('script');    
            component.setAttribute('src',"https://unpkg.com/aframe-look-at-component@0.8.0/dist/aframe-look-at-component.min.js");
            document.head.appendChild(component);

            let arnft = document.createElement('script');    
            arnft.setAttribute('src',"https://raw.githack.com/AR-js-org/AR.js/master/aframe/build/aframe-ar-nft.js");
            document.head.appendChild(arnft);

            Vue.component('a-scene', ascene)
            console.log("created")
        }
    }
</script>