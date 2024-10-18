import { makeDurationString } from "../components/ui/Statusbar/Statusbar";

// Testing the string displayed on the statusbar
describe('testStatusbar', () => {
    test('function should return a string', () => {
        const typeResult = makeDurationString(new Date(2005, 0, 24));
        expect(typeof typeResult).toBe('string');
    });

    test('function should not return a string that contains undefined', () => {
        let count = 0;
        while (count < 100) {
            const result = makeDurationString(randomDate);
            expect(result).not.toContain("undefined");
            count++;
        }
    });
});

const randomDate = getRandomDate(Date.now()-30, Date.now(), 0, 24);

function getRandomDate(start, end, startHour, endHour): Date {
    var date = new Date(+start + Math.random() * (end - start));
    var hour = startHour + Math.random() * (endHour - startHour) | 0;
    date.setHours(hour);
    return date;
  }